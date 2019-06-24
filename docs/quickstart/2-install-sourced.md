# Install source{d} Community Edition

Download the **[latest release](https://github.com/src-d/sourced-ce/releases/latest)** for your Linux, macOS (Darwin) or Windows.

## On Linux or macOS

Extract `sourced` binary from the release you downloaded, and move it into your bin folder to make it executable from any directory:

```bash
$ tar -xvf path/to/sourced-ce_REPLACE-VERSION_REPLACE-OS_amd64.tar.gz
$ sudo mv path/to/sourced-ce_REPLACE-OS_amd64/sourced /usr/local/bin/
```

## On Windows

*Please note that from now on we assume that the commands are executed in `powershell` and not in `cmd`.*

Create a directory for `sourced.exe` and add it to your `$PATH`, running these commands in a powershell as administrator:
```powershell
mkdir 'C:\Program Files\sourced'
# Add the directory to the `%path%` to make it available from anywhere
setx /M PATH "$($env:path);C:\Program Files\sourced"
# Now open a new powershell to apply the changes
```

Extract the `sourced.exe` executable from the release you downloaded, and copy it into the directory you created in the previous step:
```powershell
mv \path\to\sourced-ce_windows_amd64\sourced.exe 'C:\Program Files\sourced'
```
