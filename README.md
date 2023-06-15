# cdktf-playground repo

## follow along notes

* Installation: [docs](https://developer.hashicorp.com/terraform/tutorials/cdktf/cdktf-install)
```
npm install --global cdktf-cli@latest
cdktf --help
```

* Project setup

```
mkdir hello-world && cd hello-world
cdktf init --template=go --providers=hashicorp/null --local
```

* add provider to cdktf.json:
```
  "terraformProviders": ["null"],
```
* import provider and resource in main.go:
```

```

* General flow:
  * go-code -> synthesize to tf -> apply

aws-module

- first, add module to cdktf.json:

```
  "terraformModules": [
    {
      "name": "vpc",
      "source": "terraform-aws-modules/vpc/aws",
      "version": "3.14.2"
    }
  ],
```
- then, generate bindings:

```
cdktf get

```


