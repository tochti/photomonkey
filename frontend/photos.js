var app = angular.module('photomonkey', ['ngWebSocket'])

app.ServiceAddr = '127.0.0.1:8080';
app.ShowTime = 8 * 1000;

app.controller('photos', function($http, $log, $websocket, $scope) {
  $log.info('Hello sweet monkey!');

  $scope.photos = [];
  $scope.photo = {
    id: 'default.jpg',
    caption: 'So cute!'
  };
  $scope.photoShow = new PhotoShow($scope, app.ShowTime);

  var receiveNewPhoto = function(resp) {
    $log.info('Banana phone');
    var photo = JSON.parse(resp.data);
    $scope.photos.push(photo);

    $scope.photoShow.newPhoto();
  }

  // Init data
  var ws = $websocket('ws://' + app.ServiceAddr + '/v1/new_photos');
  ws.onMessage(receiveNewPhoto);
  ws.onError(function(event) {
    $log.error(event);
  });

  $http.get('http://' + app.ServiceAddr + '/v1/photos')
    .then(function(resp) {
      if (resp.data.length === 0) {
        return
      }

      resp.data.forEach(function(val, key) {
        $scope.photos.push(val);
      });

      $scope.photoShow.start();
    })
    .catch(function(e) {
      $log.error(e);
    })



});

var PhotoShow = function($scope, time) {
  var that = this;
  that.pos = 0;
  that.interval = 0;

  that.start = function() {
    $scope.photo = $scope.photos[that.pos];
    that.slider();
  };

  that.slider = function() {
    that.interval = setInterval(function() {
      that.pos++;
      if (that.pos > ($scope.photos.length - 1)) {
        that.pos = 0;
      }
      $scope.photo = $scope.photos[that.pos];
      $scope.$digest();
    }, time);
  };

  that.newPhoto = function() {
    clearInterval(that.interval);
    $scope.photo = $scope.photos[$scope.photos.length - 1];
    $scope.$digest();

    setTimeout(function() {
      that.slider()
    }, time);
  };
};
