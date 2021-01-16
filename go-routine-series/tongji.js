(function(window, document) {
  function Tongji() {
    this.nowTime = "";
    this.url = "";
    this.refer = ""; //refer 上一次的来源的页面 类似链表
    this.userAgent = "";
  }

  Tongji.prototype.getNowTime = function() {
    return new Date().toLocaleString();
  };
  Tongji.prototype.getUrl = function() {
    return window.location.href;
  };
  Tongji.prototype.getIp = function() {};
  Tongji.prototype.getCookie = function() {
    return document.cookie;
  };
  Tongji.prototype.getRefer = function() {
    return document.referrer;
  };
  Tongji.prototype.getUA = function() {
    return navigator.userAgent;
  };

  function random() {
    return Math.floor(Math.random() * 10000 + 500);
  }

  function formatParams(data) {
    var arr = [];
    for (var key in data) {
      arr.push(encodeURIComponent(key) + "=" + encodeURIComponent(data[key]));
    }
    arr.push("v=" + random());
    return arr.join("&");
  }
  function get(obj) {
    var url = obj.url;
    if (!url) throw new Error("must have a valid url ");
    var data = obj.data;
    if (data) {
      url += "?" + formatParams(data);
    }
    console.log(url);
    var xhr = new XMLHttpRequest();

    xhr.open("GET", url);
    xhr.onreadystatechange = function() {
      if (xhr.readyState === 4 && xhr.status == 200) {
        console.log(xhr.responseText);
      }
    };
    xhr.send();
  }

  function main() {
    var tongji = new Tongji();
    get({
      url: "xx",
      data: {
        time: tongji.getNowTime(),
        ip: tongji.getIp(),
        url: tongji.getUrl(),
        refer: tongji.getRefer()
      }
    });
  }
  window.onload = main;
})(window, document);
