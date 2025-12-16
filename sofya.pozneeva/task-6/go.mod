module task-6

go 1.22.7

replace (
    github.com/mdlayher/wifi => ./local_wifi_stub
    golang.org/x/crypto => golang.org/x/crypto v0.14.0
    golang.org/x/net => golang.org/x/net v0.17.0
    golang.org/x/sys => golang.org/x/sys v0.13.0
)

require (
    github.com/DATA-DOG/go-sqlmock v1.5.2
    github.com/stretchr/testify v1.8.4
    github.com/mdlayher/wifi v0.0.0
)
