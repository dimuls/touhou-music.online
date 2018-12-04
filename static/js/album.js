$(function() {
    
    $('.back a').click(function(event) {
        previousURL = document.referrer;
        currentURL = document.location.href;
        if (previousURL && currentURL.startsWith(previousURL)) {
            event.preventDefault();
            event.stopPropagation();
            history.back();
        }
    });
    
    var album = $('body').data('album');

    function Player() {
        var self = this;

        self.howl = null;
        self.diskNumber = null;
        self.track = null;

        self.play = function(diskNumner, track, onEnd) {
            if (self.howl !== null) {
                self.howl.stop()
            }
            self.howl = new Howl({
                src: [track.path],
                onend: onEnd,
                html5: true,
            });
            self.howl.play();
            self.diskNumber = diskNumner;
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

        self.togglePlayPauseIcon = function() {
            self.lastElement.toggleClass('fa-play');
            self.lastElement.toggleClass('fa-pause');
        };

        self.playPause = function(event) {
            diskNumber = $(event.target).data('disk-number');
            trackNumber = $(event.target).data('track-number');

            if (self.player.track !== null &&
                self.player.diskNumber === diskNumber &&
                self.player.track.number === trackNumber+'') {

                if (self.player.isPlaying()) {
                    self.player.pause();
                    self.togglePlayPauseIcon()
                } else {
                    self.player.unpause();
                    self.togglePlayPauseIcon()
                }

            } else {
                if (self.player.isPlaying()) {
                    self.togglePlayPauseIcon()
                }

                self.player.play(diskNumber,
                    self.disks[diskNumber-1].tracks[trackNumber-1],
                    function() {
                        self.togglePlayPauseIcon();
                        if (self.disks[diskNumber-1].tracks.length === trackNumber
                            && self.disks.length === diskNumber) {
                            return
                        }
                        var next = self.lastElement.parent().next();
                        while (!next.hasClass('track')) {
                            next = next.next()
                        }
                        next.children('.play').click()
                    }
                );
                self.lastElement = $(event.target);
                self.togglePlayPauseIcon();
            }

        };

        $('.play').click(self.playPause);
    }

    window.app = new Album(album);
});