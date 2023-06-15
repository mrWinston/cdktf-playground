# Ec2 Instance Example

This example deploys an EC2 instance into a newly created VPC. It shows handling of the following scenarios:

- using a module from the terraform registry
- using module outputs as inputs to other resources
- Using `data`-resources

## How to run

Run the following commands in this directory to apply the stack:
```
# cdktf get generates the go language bindings for external modules dynamically
cdktf get 
cdktf apply
```

