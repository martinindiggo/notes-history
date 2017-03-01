// myApp is a sample AngularJS app.
var myApp = angular.module("myApp", []);

// notesController is a sample controller to implement basic notes functionality.
angular.module("myApp").controller("notesController", ["$scope", "$timeout", "notesFactory", ($scope, $timeout, notesFactory) => {
    
    let meetingId = "1";

    $scope.currentNote = {};
    $scope.allNotes = [];

    $scope.addNote = (note) => {
        
        let newNote = {};
        newNote.author = note.author;
        newNote.meetingId = meetingId;
        newNote.content = note.content;
        notesFactory.addNote(newNote).then((res) => {
            $timeout(() => {
                init();
                $scope.noteForm.$setPristine();
            }, 500);
        });

    }

    $scope.deleteNote = (id) => {
        if (id) {
            notesFactory.deleteNote(id).then((res) => {
                $timeout(() => {
                    init();
                    $scope.noteForm.$setPristine();
                }, 500);
            });
        }
    }

    function init() {
        notesFactory.getAllNotes(meetingId).then((res) => {
            $scope.allNotes = res.data;
        });

        notesFactory.getLatestNote(meetingId).then((res) => {
            if (res.data && res.data.length > 0) {
                $scope.currentNote = res.data[0];  
            }
        });
    }

    init();

}]);

// notesFactory is a sample factory for working with notes.
angular.module("myApp").factory("notesFactory", ["$http", ($http) => {
    var factory = {};

    factory.getAllNotes = (meetingId) => {
        return $http.get("api/notes/meeting/" + meetingId);
    }

    factory.getLatestNote = (meetingId) => {
        return $http.get(`api/notes/meeting/${meetingId}/latest`);
    }

    factory.addNote = (note) => {
        return $http.post("api/notes", note);
    }

    factory.deleteNote = (id) => {
        return $http.delete("api/notes/" + id);
    }

    return factory;
}]);