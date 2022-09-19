# k8scp

A command line tool for copying files to K8s pods

## svr
```
./svr -k ~/kubeconfig
```

## scp
```
./scp -u http://127.0.0.1:8080/upload -s ~/test.txt -n test-ns -p nginx-0 -c nginx -d '/data/'
```