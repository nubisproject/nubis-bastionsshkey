## nubis-bastionsshkey
Attempts to query ldap and insert it into consul

### Requirements
1. Make sure `GOPATH` is set
2. Requires the following package

    ```
    go get -u go.mozilla.org/mozldap
    go get -u gopkg.in/yaml.v2
    go get -u github.com/hashicorp/consul/api
    go get -u github.com/aws/aws-sdk-go
    ```

3. Run script for testing by doing this:

    ```
    go run *.got
    ```
