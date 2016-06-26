var app = angular.module('photomonkey', [
  'ngRoute',
  'photomonkey.ctrls',
]);

app.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/default', {
    templateUrl: '/public/tpls/main.html',
    controller: 'PhotosCtrl',
  }).otherwise('/default');
}
]);
