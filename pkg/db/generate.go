package db

//go:generate sh -c "rm -rf mocks && mkdir -p mock"
//go:generate ../../../bin/minimock -i TxManager -o ./mock -s "_minimock.go"
