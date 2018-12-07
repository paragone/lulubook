echo "set GOPATH"
export GOPATH=`pwd`

if [ "$1" == "china" ];then
	echo "for china,  get golang.org/x/net"

	git clone https://github.com/golang/net.git

	mkdir -p src/golang.org/x/

	mv net src/golang.org/x/
fi

echo "go get gin"
go get  github.com/gin-gonic/gin

echo "go get goquery"

go get github.com/PuerkitoBio/goquery

echo "go get string convert"

go get github.com/axgle/mahonia

echo "go get mongodb driver"

go get github.com/globalsign/mgo

