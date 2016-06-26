var ctrls = angular.module('photomonkey.ctrls', ['ngWebSocket']);

ctrls.controller('PhotosCtrl', function($http, $log, $websocket, $scope) {
  $log.info('Hello sweet monkey!');

  var serviceAddr = '127.0.0.1:8080';
  var showTime = 2 * 1000;

  $scope.photos = [];
  $scope.photo = {};

  var slider = {
    pos: 0,
    interval: 0,

    start: function() {
      this.setPhoto(this.pos);
      this.slide();
    },

    slide: function() {
      var that = this;
      that.interval = setInterval(function() {
        console.log('slider');
        that.pos++;
        if (that.pos > ($scope.photos.length - 1)) {
          that.pos = 0;
        }
        that.setPhoto(that.pos);
      }, showTime);
    },

    newPhoto: function() {
      clearInterval(this.interval);
      that.setPhoto($scope.photos.length - 1);

      setTimeout(function() {
        this.slide()
      }, showTime);
    },

    setPhoto: function(pos) {
      $scope.photo = $scope.photos[pos];
      $('#photo').attr('src', '/files/' + $scope.photo.id);
      $('#photo').attr('alt', $scope.photo.caption);
      console.log($scope.photo);
    },
  };

  var receiveNewPhoto = function(resp) {
    $log.info('Banana phone');
    var photo = JSON.parse(resp.data);
    $scope.photos.push(photo);

    slider.newPhoto();
  }

  // Init data
  $('#photo').on('load', function() {
    console.log('Photo is loaded');
    var windowHeight = window.innerHeight;
    var windowWidth = window.innerWidth;

    var imgHeight = $('#photo').height();
    var imgWidth = $('#photo').width();

    var newImgHeight = 0;
    var newImgWidth = 0;

    if (imgHeight > imgWidth) {
      if (imgHeight > windowHeight) {
        newImgHeight = windowHeight * 0.95;
      } else {
        newImgHeight = imgHeight;
      }

      var f = (1 / imgHeight) * newImgHeight;
      newImgWidth = imgWidth * f;

    } else {
      if (imgWidth > windowWidth) {
        newImgWidth = windowWidth * 0.95;
      } else {
        newImgWidth = imgWidth;
      }

      var f = (1 / imgWidth) * newImgWidth;
      newImgHeight = imgHeight * f;
    }

    console.log('window Height', windowHeight);
    console.log('window Width', windowWidth);
    console.log('ImgHeight', imgHeight);
    console.log('ImgWidth', imgWidth);
    console.log('newImgHeight', newImgHeight);
    console.log('newImgWidth', newImgWidth);

    //$(this).css('height', newImgHeight);
    //$(this).css('width', newImgWidth);
    $(this).height(newImgHeight);
    $(this).width(newImgWidth);

  });

  var ws = $websocket('ws://' + serviceAddr + '/v1/new_photos');
  ws.onMessage(receiveNewPhoto);
  ws.onError(function(event) {
    $log.error(event);
  });

  $http.get('http://' + serviceAddr + '/v1/photos')
    .then(function(resp) {
      if (resp.data.length === 0) {
        return
      }

      resp.data.forEach(function(val, key) {
        $scope.photos.push(val);
      });

      console.log($scope.photos);
      slider.start();
    })
    .catch(function(e) {
      $log.error(e);
    })

});
