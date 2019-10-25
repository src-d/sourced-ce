package workdir

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/pbnjay/memory"
	"github.com/pkg/errors"
	"github.com/serenize/snaker"
	composefile "github.com/src-d/sourced-ce/cmd/sourced/compose/file"
)

// InitLocal initializes the workdir for local path and returns the Workdir instance
func InitLocal(reposdir string) (*Workdir, error) {
	dirName := encodeDirName(reposdir)
	envf := newLocalEnvFile(dirName, reposdir)

	return initialize(dirName, "local", envf)
}

// InitOrgs initializes the workdir for organizations and returns the Workdir instance
func InitOrgs(orgs []string, token string, withForks bool) (*Workdir, error) {
	// be indifferent to the order of passed organizations
	sort.Strings(orgs)
	dirName := encodeDirName(strings.Join(orgs, ","))

	envf := envFile{}
	err := readEnvFile(dirName, "orgs", &envf)
	if err == nil && envf.NoForks == withForks {
		return nil, ErrInitFailed.Wrap(
			fmt.Errorf("workdir was previously initialized with a different value for forks support"))
	}
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	// re-create env file to make sure all fields are updated
	envf = newOrgEnvFile(dirName, orgs, token, withForks)

	return initialize(dirName, "orgs", envf)
}

func readEnvFile(dirName string, subPath string, envf *envFile) error {
	workdir, err := workdirPath(dirName, subPath)
	if err != nil {
		return err
	}

	envPath := filepath.Join(workdir, ".env")
	b, err := ioutil.ReadFile(envPath)
	if err != nil {
		return err
	}

	return envf.UnmarshalEnv(b)
}

func encodeDirName(dirName string) string {
	return base64.URLEncoding.EncodeToString([]byte(dirName))
}

func workdirPath(dirName string, subPath string) (string, error) {
	path, err := workdirsPath()
	if err != nil {
		return "", err
	}

	workdir := filepath.Join(path, subPath, dirName)
	if err != nil {
		return "", err
	}

	return workdir, nil
}

func initialize(dirName string, subPath string, envf envFile) (*Workdir, error) {
	path, err := workdirsPath()
	if err != nil {
		return nil, err
	}

	workdir, err := workdirPath(dirName, subPath)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(workdir, 0755)
	if err != nil {
		return nil, errors.Wrap(err, "could not create working directory")
	}

	defaultFilePath, err := composefile.InitDefault()
	if err != nil {
		return nil, err
	}

	composePath := filepath.Join(workdir, "docker-compose.yml")
	if err := link(defaultFilePath, composePath); err != nil {
		return nil, err
	}

	envPath := filepath.Join(workdir, ".env")
	contents, err := envf.MarshalEnv()
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(envPath, contents, 0644)

	if err != nil {
		return nil, errors.Wrap(err, "could not write .env file")
	}

	b := &builder{workdirsPath: path}
	return b.Build(workdir)
}

type envFile struct {
	ComposeProjectName string

	GitbaseVolumeType   string
	GitbaseVolumeSource string
	GitbaseSiva         bool

	GithubOrganizations []string
	GithubToken         string

	NoForks bool

	GitbaseLimitCPU      float32
	GitcollectorLimitCPU float32
	GitbaseLimitMem      uint64
}

func newLocalEnvFile(dirName, repoDir string) envFile {
	f := envFile{
		ComposeProjectName: fmt.Sprintf("srcd-%s", dirName),

		GitbaseVolumeType:   "bind",
		GitbaseVolumeSource: repoDir,
	}
	f.addResourceLimits()

	return f
}

func newOrgEnvFile(dirName string, orgs []string, token string, withForks bool) envFile {
	f := envFile{
		ComposeProjectName: fmt.Sprintf("srcd-%s", dirName),

		GitbaseVolumeType:   "volume",
		GitbaseVolumeSource: "gitbase_repositories",
		GitbaseSiva:         true,

		GithubOrganizations: orgs,
		GithubToken:         token,

		NoForks: !withForks,
	}
	f.addResourceLimits()

	return f
}

func (f *envFile) addResourceLimits() {
	// limit CPU for containers
	dockerCPUs, err := dockerNumCPU()
	if err != nil { // show warning
		fmt.Println(err)
	}
	// apply gitbase resource limits only when docker runs without any global limits
	// it's default behaviour on linux
	if runtime.NumCPU() == dockerCPUs {
		f.GitbaseLimitCPU = float32(dockerCPUs) - 0.1
	}
	// always apply gitcollector limit
	if dockerCPUs > 0 {
		halfCPUs := float32(dockerCPUs) / 2.0
		// let container consume more than a half if there is only one cpu available
		// otherwise it will be too slow
		if halfCPUs < 1 {
			halfCPUs = 1
		}
		f.GitcollectorLimitCPU = halfCPUs - 0.1
	}

	// limit memory for containers
	dockerMem, err := dockerTotalMem()
	if err != nil { // show warning
		fmt.Println(err)
	}
	// apply memory limits only when only when docker runs without any global limits
	// it's default behaviour on linux
	if dockerMem == memory.TotalMemory() {
		f.GitbaseLimitMem = uint64(float64(dockerMem) * 0.9)
	}
}

var newlineChar = "\n"

func init() {
	if runtime.GOOS == "windows" {
		newlineChar = "\r\n"
	}
}

// implementation can be moved to separate package if we need to marshal any other structs
// supports only simple types
func (f envFile) MarshalEnv() ([]byte, error) {
	var b bytes.Buffer

	v := reflect.ValueOf(f)
	rType := v.Type()

	for i := 0; i < rType.NumField(); i++ {
		field := rType.Field(i)
		fieldEl := v.Field(i)
		if field.Anonymous {
			panic("struct composition isn't supported")
		}

		name := strings.ToUpper(snaker.CamelToSnake(field.Name))
		switch field.Type.Kind() {
		case reflect.Slice:
			slice := make([]string, fieldEl.Len())
			for i := 0; i < fieldEl.Len(); i++ {
				slice[i] = fmt.Sprintf("%v", fieldEl.Index(i).Interface())
			}
			fmt.Fprintf(&b, "%s=%v%s", name, strings.Join(slice, ","), newlineChar)
		case reflect.Bool:
			// marshal false value as empty string instead of "false" string
			if fieldEl.Interface().(bool) {
				fmt.Fprintf(&b, "%s=true%s", name, newlineChar)
			} else {
				fmt.Fprintf(&b, "%s=%s", name, newlineChar)
			}
		default:
			fmt.Fprintf(&b, "%s=%v%s", name, fieldEl.Interface(), newlineChar)
		}
	}

	return b.Bytes(), nil
}

// implementation can be moved to separate package if we need to unmarshal any other structs
// supports only simple types
func (f *envFile) UnmarshalEnv(b []byte) error {
	v := reflect.ValueOf(f).Elem()

	r := bytes.NewReader(b)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || !strings.Contains(line, "=") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		name := parts[0]
		value := parts[1]
		field := v.FieldByName(snaker.SnakeToCamel(strings.ToLower(name)))
		// skip unknown values
		if !field.IsValid() {
			continue
		}
		// skip empty values
		if value == "" {
			continue
		}
		switch field.Kind() {
		case reflect.String:
			field.SetString(value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("can't parse variable %s with value %s: %v", name, value, err)
			}
			field.SetInt(i)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			i, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return fmt.Errorf("can't parse variable %s with value %s: %v", name, value, err)
			}
			field.SetUint(i)
		case reflect.Float32, reflect.Float64:
			i, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("can't parse variable %s with value %s: %v", name, value, err)
			}
			field.SetFloat(i)
		case reflect.Bool:
			if value == "true" {
				field.SetBool(true)
			} else {
				field.SetBool(false)
			}
		case reflect.Slice:
			if field.Type().Elem().Kind() != reflect.String {
				panic("only slices of strings are supported")
			}
			vs := strings.Split(value, ",")
			slice := reflect.MakeSlice(field.Type(), len(vs), len(vs))
			for i, v := range vs {
				slice.Index(i).SetString(v)
			}
			field.Set(slice)
		default:
			panic(fmt.Sprintf("unsupported type: %v", field.Kind()))
		}
	}

	return scanner.Err()
}

// returns number of CPUs available to docker
func dockerNumCPU() (int, error) {
	// use cli instead of connection to docker server directly
	// in case server exposed by http or non-default socket path
	info, err := exec.Command("docker", "info", "--format", "{{.NCPU}}").Output()
	if err != nil {
		return 0, err
	}

	cpus, err := strconv.Atoi(strings.TrimSpace(string(info)))
	if err != nil || cpus == 0 {
		return 0, fmt.Errorf("Couldn't get number of available CPUs in docker")
	}

	return cpus, nil
}

// returns total memory in bytes available to docker
func dockerTotalMem() (uint64, error) {
	info, err := exec.Command("docker", "info", "--format", "{{.MemTotal}}").Output()
	if err != nil {
		return 0, err
	}

	mem, err := strconv.ParseUint(strings.TrimSpace(string(info)), 10, 64)
	if err != nil || mem == 0 {
		return 0, fmt.Errorf("Couldn't get of available memory in docker")
	}

	return mem, nil
}
