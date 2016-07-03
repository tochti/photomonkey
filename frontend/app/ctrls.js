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
      this.setPhoto($scope.photos.length - 1);

      var that = this;
      setTimeout(function() {
        that.slide()
      }, showTime);
    },

    setPhoto: function(pos) {
      console.log('--- set photo ---');

      $scope.photo = $scope.photos[pos];
      photoPath = '/files/' + $scope.photo.id;

      var img = new Image();

      img.addEventListener('load', function() {
        var windowHeight = window.innerHeight * 0.95;
        var windowWidth = window.innerWidth * 0.95;

        var imgHeight = img.height;
        var imgWidth = img.width;

        var newImgHeight = 0;
        var newImgWidth = 0;

        if (imgHeight > windowHeight && imgWidth > windowWidth) {
          console.log('width and height to big');
          // Wenn das Bild auf beiden Seiten zu groß ist
          // finde heraus über welche Seite das Bild weiter heraus steht.
          // Passe diese Seite auf Fenstergröße an und verwende den Faktor der Veränderung
          // um die andere Seite damit anzupassen da diese kleiner war wird
          // diese daher auch kleiner sein als das Fenster.
          dH = imgHeight - windowHeight;
          dW = imgWidth - windowWidth;

          sizes = {}
          if (dH > dW) {
            newImgHeight = windowHeight;
            var f = (1 / imgHeight) * newImgHeight;
            console.log("factor ", f);
            newImgWidth = imgWidth * f;
          } else {
            newImgWidth = windowWidth;
            var f = (1 / imgWidth) * newImgWidth;
            console.log("factor ", f);
            newImgHeight = imgHeight * f;
          }

        } else if (imgHeight > windowHeight) {
          // Wenn nur eine Seite des Bilds über das Fenster hinausragt
          // wird diese Seite auf Fenster größe angepasst die ander Seite wird um den Faktor
          // der Veränderung angepasst.
          console.log('height to big');
          if (imgHeight > windowHeight) {
            newImgHeight = windowHeight;
          } else {
            newImgHeight = imgHeight;
          }

          var f = (1 / imgHeight) * newImgHeight;
          console.log("factor ", f);
          newImgWidth = imgWidth * f;

        } else if (imgWidth > windowWidth) {
          console.log('width to big');
          if (imgWidth > windowWidth) {
            newImgWidth = windowWidth;
          } else {
            newImgWidth = imgWidth;
          }

          var f = (1 / imgWidth) * newImgWidth;
          console.log("factor ", f);
          newImgHeight = imgHeight * f;
        } else {
          // Keine Seite des Bilds ist zu groß
          newImgHeight = imgHeight;
          newImgWidth = imgWidth;
        }

        //console.log(photoPath);
        console.log(img);
        console.log('window Height', windowHeight);
        console.log('window Width', windowWidth);
        console.log('ImgHeight', imgHeight);
        console.log('ImgWidth', imgWidth);
        console.log('newImgHeight', newImgHeight);
        console.log('newImgWidth', newImgWidth);

        photoItem = $('#photo');
        photoItem.attr('src', photoPath);
        photoItem.attr('alt', $scope.photo.caption);
        photoItem.width(newImgWidth);
        photoItem.height(newImgHeight);
        if ($scope.photo.caption.length !== 0) {
          $('#caption').html($scope.photo.caption);
          $('#caption').show();
        } else {
          $('#caption').hide();
        }
        box = $('#box');
        box.width(newImgWidth);
        box.height(newImgHeight);
        box.css('margin-top', (windowHeight - newImgHeight) / 2);

      });

      img.src = photoPath;

    },
  };

  var receiveNewPhoto = function(resp) {
    $log.info('Banana phone');
    var photo = JSON.parse(resp.data);
    $scope.photos.push(photo);

    slider.newPhoto();
  }

  // Init data
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
