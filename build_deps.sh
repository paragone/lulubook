echo "set GOPATH"
export GOPATH=`pwd`

echo "go get gin"
go get  github.com/gin-gonic/gin

echo "for china,  get golang.org/x/net"

git clone https://github.com/golang/net.git

mkdir -p golang.org/x/

mv net golang.org/x/

echo "go get goquery"

go get github.com/PuerkitoBio/goquery

echo "go get string convert"

go get github.com/axgle/mahonia

echo "go get mongodb driver"

go get github.com/globalsign/mgo

