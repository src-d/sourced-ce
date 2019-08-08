:: we can't use makefile for windows because it depends on CI makefile which depends on shell

:: compile sourced-ce
go build -tags=forceposix -o build/sourced-ce_windows_amd64/sourced.exe ./cmd/sourced

:: run tests
set TMPDIR_INTEGRATION_TEST=C:\tmp
:: see https://stackoverflow.com/questions/22948189/batch-getting-the-directory-is-not-empty-on-rmdir-command
del /f /s /q %TMPDIR_INTEGRATION_TEST% 1>nul
rmdir /s /q %TMPDIR_INTEGRATION_TEST%
mkdir %TMPDIR_INTEGRATION_TEST%
set TMP=%TMPDIR_INTEGRATION_TEST%
go test -timeout 20m -parallel 1 -count 1 -tags="forceposix integration" github.com/src-d/sourced-ce/test/ -v
