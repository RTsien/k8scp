# k8scp

A command line tool for copying files to K8s pods

## svr
```
./svr -k ~/kubeconfig -p 8080
```

## scp
```
# Linux/macOS
./scp -u http://127.0.0.1:8080/upload -s ~/test.txt -n test-ns -p nginx-0 -c nginx -d '/data/'

# Win CMD/PowerShell
./scp -u http://127.0.0.1:8080/upload -s C:\test.txt -n test-ns -p nginx-0 -c nginx -d '/data/'

# Win Git-Bash ref: https://stackoverflow.com/q/7250130
./scp -u http://127.0.0.1:8080/upload -s ~/test.txt -n test-ns -p nginx-0 -c nginx -d '//data/'
```

```
# Support uploading the entire directory
./scp -u http://127.0.0.1:8080/upload -s ~/ -n test-ns -p nginx-0 -c nginx -d '/data/'
```

```
Usage:
  scp [flags]

Flags:
  -c, --container string   container name
  -d, --dst string         destination file path
  -h, --help               help for scp
  -n, --namespace string   k8s namespace
  -p, --pod string         pod name
  -r, --rate int           rate limit MiB/s (default 12)
  -s, --src string         source file path
  -u, --url string         server url
```
