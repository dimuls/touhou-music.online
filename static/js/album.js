$(function() {
    var album = $("body").data("album");

    function Player() {
        var self = this;

        self.howl = null;
        self.track = null;

        self.play = function(track, onEnd) {
            if (self.howl !== null) {
                self.howl.stop()
            }
            self.howl = new Howl({
                src: [track.path],
                onend: onEnd
            });
            self.howl.play();
            self.track = track;
        };

        self.stop = function() {
            if (self.howl !== null) {
                self.howl.stop();
                self.howl = null;
                self.track = null;
            }
        };

        self.pause = function() {
            self.howl.pause();
        };

        self.unpause = function() {
            self.howl.play();
        };

        self.isPlaying = function() {
            if (self.howl === null) {
                return false;
            }
            return self.howl.playing();
        }
    }

    function Album(data) {
        var self = this;

        $.extend(self, data);

        self.player = new Player();
        self.lastElement = null;

        self.playPause = function(_, event) {
            trackNumber = $(event.target).data("track-number");

            if (self.player.track !== null &&
                self.player.track.number === trackNumber+'') {

                if (self.player.isPlaying()) {
                    self.player.pause();
                    self.lastElement.toggleClass("paused")
                } else {
                    self.player.unpause();
                    self.lastElement.toggleClass("paused")
                }

            } else {
                if (self.player.isPlaying()) {
                    self.lastElement.toggleClass("paused")
                }
                self.player.play(self.tracks[trackNumber-1], function() {
                    if (self.tracks.length === trackNumber) {
                        self.lastElement.toggleClass("paused");
                        return
                    }
                    self.lastElement.toggleClass("paused");
                    self.lastElement.parent().next().children(".play").click()
                });
                self.lastElement = $(event.target);
                self.lastElement.toggleClass("paused")
            }

        };
    }

    ko.applyBindings(new Album(album));
});