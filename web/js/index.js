(function() {
    //html5本地存储,localstorage在同一域名下共享
    var Util = (function() {
        var prefix = 'html5_reader_';
        var StorageGetter = function(key) {
            return localStorage.getItem(prefix + key);
        };
        var StorageSetter = function(key, val) {
            return localStorage.setItem(prefix + key, val);
        };
        return {
            StorageGetter: StorageGetter,
            StorageSetter: StorageSetter
        }
    })();

    var Dom = {
        top_nav: $('#top_nav'),
        bottom_nav: $('#bottom_nav'),
        font_container: $('.font-container'),
        font_button: $('#font-button'),
        bk_container: $('.bk-container'),
    };
    var mode = 0;
    var Win = $(window);
    var Doc = $(document);
    var readModel;
    var readBookList;
    var readChapterList;
    var readChapterContent;
    var RootContainer = $('#fiction_container');
    var initFontSize = Util.StorageGetter('font_size');
    var backgroundColor = Util.StorageGetter('background');
    initFontSize = parseInt(initFontSize);
    if (!initFontSize) {
        initFontSize = 14;
    }
    RootContainer.css('font-size', initFontSize);
    RootContainer.css('background', backgroundColor);

    function main() {
        //整个项目的入口函数
        readModel = ReadModel();
        readBookList = RenderBookList(RootContainer);
        readChapterList = RenderChapterList(RootContainer);
        readChapterContent= RenderChapterContent(RootContainer);
        if(Util.StorageGetter("book") == "undefined" || Util.StorageGetter("chapter") == "undefined"
        ||Util.StorageGetter("book") == null || Util.StorageGetter("chapter") == null) {
            readModel.init(function(data) {
                readBookList(data);
            });
        } else {
            //var bookid = parseInt(Util.StorageGetter("book"));
            var chapterid = parseInt(Util.StorageGetter("chapter"));
            //readModel.listChapter(bookid, readChapterList);
            readModel.curChapter(chapterid, readChapterContent);

        }


        EventHanlder();
    }
    //数据层
    function ReadModel() {
        //实现和阅读器先关的数据交互的方法
        var Chapter_id;
        var Book_id;

        var limited = 20;
        var offset = 0;
        //初始化函数
        var init = function(UIcallback) {
            getBookList(function(data) {
                UIcallback && UIcallback(data);
            });
        };

        //获取书籍列表
        var getBookList = function(callback) {
            if(mode != 0){
                limited = 20;
                offset = 0;
            }
            $.getJSON('/api/v1/view/?offset='+offset+'&limited='+limited, function(data) {
                mode = 0;
                callback && callback(data);
            })
        };
        //获取书籍章节
        var getChapterOfBook = function(book_id, callback) {
            if(mode != 1){
                limited = 20;
                offset = 0;
            }
            $.getJSON('/api/v1/view/'+book_id+'/?offset='+offset+'&limited='+limited, function(data) {
                Book_id = book_id;
                mode = 1;

                Util.StorageSetter('book', Book_id);
                callback && callback(data);
            });

        };

        //获取章节内容
        var getCurChapterContent = function(chapter_id, callback) {
            if(mode != 2){
                limited = 20;
                offset = 0;
            }

            if(Book_id == undefined)
                Book_id = Util.StorageGetter('book');

            $.getJSON('/api/v1/view/'+Book_id+'/'+chapter_id+'/', function(data) {
                Chapter_id = chapter_id;
                mode = 2;

                Util.StorageSetter('chapter', Chapter_id);
                callback && callback(data);
            });

        };

        var listBook = function(UIcallback) {
            getBookList(function(data) {
                UIcallback && UIcallback(data);
            });
        };
        var listChapter = function(book_id,UIcallback) {

            getChapterOfBook(book_id,function(data) {
                UIcallback && UIcallback(data);
            });
        };
        var curChapter = function(chapter_id, UIcallback) {

            getCurChapterContent(chapter_id,function(data) {
                UIcallback && UIcallback(data);
            });
        };
        //上一章事件处理函数
        var prevBtn = function(UIcallback) {
            switch (mode){
                case 0:{
                    if (offset === 0) {
                        return;
                    }
                    offset -= 20;
                    getBookList(UIcallback);
                    break;
                }
                case 1:{
                    if (offset === 0) {
                        return;
                    }
                    offset -= 20;
                    getCurChapterContent(Book_id, UIcallback);
                    break;
                }
                case 2:{
                    if (Chapter_id === 0) {
                        return;
                    }
                    Chapter_id -= 1;
                    getCurChapterContent(Chapter_id, UIcallback);
                    break;
                }
                default:
                    return;
            }


        };
        //下一章事件处理函数
        var nextBtn = function(UIcallback) {
            switch (mode) {
                case 0: {
                    offset += 20;
                    getBookList(UIcallback);
                    break;
                }
                case 1: {
                    offset += 20;
                    getCurChapterContent(Book_id, UIcallback);
                    break;
                }
                case 2: {
                    Chapter_id += 1;
                    getCurChapterContent(Chapter_id, UIcallback);
                    break;
                }
                default:
                    return;
            }
        };
        return {
            init: init,
            listBook:listBook,
            listChapter: listChapter,
            curChapter,curChapter,
            prevBtn: prevBtn,
            nextBtn: nextBtn,
        }
    }

    function RenderBookList(container) {

        mode = 0;
        //todo 渲染基本的UI结构
        function parseBookList(jsonData) {
            var html = '<h1>书籍列表：</h1>';
            for (var i = 0; i < jsonData.length; i++) {
                html += '<div class="booklist" id="'+jsonData[i]._id+'"><h4>'
                    + jsonData[i]._id + "、" +jsonData[i].name + '</h4></div>';
            }
            return html;
        }
        return function(data) {
            container.html(parseBookList(data));
            $('#action_mid').hide();
            Dom.top_nav.hide();
            Dom.bottom_nav.hide();
            Dom.font_container.hide();
            $('.icon-ft').removeClass('current');
        }
    }
    function RenderChapterList(container) {

        mode = 1;
        //todo 渲染基本的UI结构
        function parseChapterList(jsonData) {
            var html = '<h1>目录：</h1>';
            for (var i = 0; i < jsonData.length; i++) {
                html += '<div class="chapterlist" id="'+jsonData[i]._id+'"><h4>'+ jsonData[i].title + '</h4></div>';
            }
            return html;
        }
        return function(data) {
            container.html(parseChapterList(data));
            $('#action_mid').hide();
            Dom.top_nav.hide();
            Dom.bottom_nav.hide();
            Dom.font_container.hide();
            $('.icon-ft').removeClass('current');
        }
    }
    function RenderChapterContent(container) {

        mode = 2;
        //todo 渲染基本的UI结构
        function parseChapterContent(jsonData) {
            var html = '<h1>'+ jsonData.title +'</h1>';
            html += '<p>'+ jsonData.content +'</p>';
            return html;
        }
        return function(data) {
            container.html(parseChapterContent(data));
            $('#action_mid').show();
            Dom.top_nav.hide();
            Dom.bottom_nav.hide();
            Dom.font_container.hide();
            $('.icon-ft').removeClass('current');
        }
    }

    function pageScroll(){

        //获取scrollTop值
        var sTop=window.pageYOffset||document.documentElement.scrollTop|| document.body.scrollTop|| 0;
        //把内容滚动指定的像素数（第一个参数向右滚动的像素数，第二个参数向下滚动的像素数）
        window.scrollBy(0,-sTop);

    }

    function EventHanlder() {
        //控制层的作用
        //交互的事件绑定
        //安卓4.0前，click事件有一定延迟（300ms）
        //zepto 模拟的点击tap事件
        //控制顶部和底部导航栏的显示与隐藏
        $('#action_mid').click(function() {
            if (Dom.top_nav.css('display') == "none") {
                Dom.top_nav.show();
                Dom.bottom_nav.show();
            } else {
                Dom.top_nav.hide();
                Dom.bottom_nav.hide();
                Dom.font_container.hide();
                $('.icon-ft').removeClass('current');
            }

        });
        Dom.font_button.click(function() {
            if (Dom.font_container.css('display') == "none") {
                Dom.font_container.show();
                $('.icon-ft').addClass('current');

            } else {
                Dom.font_container.hide();
                $('.icon-ft').removeClass('current');
            }
        });
        $('#large-font').click(function() {
            if (initFontSize > 20)
                return;
            initFontSize += 1;
            RootContainer.css('font-size', initFontSize);
            Util.StorageSetter('font_size', initFontSize);
        });

        $('#small-font').click(function() {
            if (initFontSize < 12)
                return;
            initFontSize -= 1;
            RootContainer.css('font-size', initFontSize);
            Util.StorageSetter('font_size', initFontSize);
        });
        //设置背景颜色，each函数用来遍历
        $.each(Dom.bk_container, function(index, value) {
            Dom.bk_container[index].onclick = function() {
                backgroundColor = $(Dom.bk_container[index]).css('background').slice(0, 18);
                RootContainer.css('background', backgroundColor);
                Util.StorageSetter('background', backgroundColor);
            };
        });
        $('#top_nav').click(function() {
            readModel.listBook(function(data) {
                readBookList(data);
            });
        });
        $('#menu_button').click(function() {
            readModel.listBook(function(data) {
                readBookList(data);
            });
        });
        $('#fiction_container').on('click','.booklist',function(){
            readModel.listChapter(parseInt(this.id),readChapterList);
        });
        $('#fiction_container').on('click','.chapterlist',function(){
            readModel.curChapter(parseInt(this.id),readChapterContent);
        });


        $('#daytime-button').click(function() {
            $('#daytime-button').hide();
            $('#night-button').show();
            RootContainer.css('background', '#e9dfc7');
            Util.StorageSetter('background', '#e9dfc7');
        });
        $('#night-button').click(function() {
            $('#daytime-button').show();
            $('#night-button').hide();
            RootContainer.css('background', '#000');
            Util.StorageSetter('background', '#000');
        });
        //鼠标滑动时触发事件处理
        Win.scroll(function() {
            Dom.top_nav.hide();
            Dom.bottom_nav.hide();
            Dom.font_container.hide();
            $('.icon-ft').removeClass('current');
        });
        //章节翻页，先获得章节的翻页数据，再把数据拿出来渲染
        $('#prev_button').click(function() {
            readModel.prevBtn(function(data) {
                switch (mode) {
                    case 0: {
                        readBookList(data);
                        break;
                    }
                    case 1: {
                        readChapterList(data);
                        break;
                    }
                    case 2: {
                        readChapterContent(data);
                        break;
                    }
                    default:
                        return;
                }
            });
            pageScroll();
        });
        $('#next_button').click(function() {
            readModel.nextBtn(function(data) {
                switch (mode) {
                    case 0: {
                        readBookList(data);
                        break;
                    }
                    case 1: {
                        readChapterList(data);
                        break;
                    }
                    case 2: {
                        readChapterContent(data);
                        break;
                    }
                    default:
                        return;
                }
            });
            pageScroll();
        });
    }
    main();
})();
