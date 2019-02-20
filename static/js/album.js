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

        self.titleTextEl = $('.player .title .text');
        self.titleTimeEl = $('.player .title .time');

        self.howl = null;
        self.discNumber = null;
        self.track = null;

        self.timeUpdater = null;

        self.updateTitleTime = function() {
            if (self.howl === null) {
                return
            }
            var total = self.howl.duration();
            var totalMinutes = Math.floor(total / 60);
            var totalSeconds = ("0" +Math.round(total - totalMinutes * 60)).slice(-2);

            var seek = self.howl.seek();
            var minutes = Math.floor(seek/ 60);
            var seconds = ("0" + Math.round(seek - minutes * 60)).slice(-2);

            var time = minutes + ":" + seconds + " / "
                + totalMinutes + ":" + totalSeconds;

            self.titleTimeEl.text(time);
        };

        self.preload = function(discNumber, track, onEnd) {
            self.titleTextEl.text(track.number + '. ' + track.title);
            self.howl = new Howl({
                src: [track.path],
                onend: onEnd,
                onplay: function() {
                    self.timeUpdater = setInterval(self.updateTitleTime, 1000);
                },
                onload: function() {
                    self.updateTitleTime();
                }
            });
            self.track = track;
            self.discNumber = discNumber;
        };

        self.play = function(discNumber, track, onEnd) {
            if (self.howl !== null) {
                if (self.discNumber === discNumber &&
                    self.track.number === track.number) {
                    self.howl.play();
                    return
                }
                self.howl.stop()
            }
            self.titleTextEl.text(track.number + '. ' + track.title);
            self.titleTimeEl.text("");
            self.howl = new Howl({
                src: [track.path],
                onend: onEnd,
                onplay: function() {
                    self.timeUpdater = setInterval(self.updateTitleTime, 1000);
                },
                onload: function() {
                    self.updateTitleTime();
                },
                html5: true,
            });
            self.howl.play();
            self.discNumber = discNumber;
            self.track = track;
        };

        self.stop = function() {
            if (self.howl !== null) {
                self.howl.stop();
                self.howl = null;
                self.track = null;
                clearInterval(self.timeUpdater);
            }
        };

        self.pause = function() {
            self.howl.pause();
            clearInterval(self.timeUpdater);
        };

        self.unpause = function() {
            self.howl.play();
            self.timeUpdater = setInterval(self.updateTitleTime, 1000);
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

        self.toggleActiveTrack = function() {
            self.lastElement.parent().toggleClass('active')
        };

        self.playPause = function(event) {
            discNumber = $(event.target).data('disc-number');
            trackNumber = $(event.target).data('track-number');

            if (self.player.track !== null &&
                self.player.discNumber === discNumber &&
                self.player.track.number === trackNumber+'') {

                if (self.player.isPlaying()) {
                    self.player.pause();
                    self.togglePlayPauseIcon();
                } else {
                    self.player.unpause();
                    self.togglePlayPauseIcon();
                }

            } else {
                if (self.player.isPlaying()) {
                    self.togglePlayPauseIcon();
                    self.toggleActiveTrack();
                }

                self.player.play(discNumber,
                    self.discs[discNumber-1].tracks[trackNumber-1],
                    function() {
                        self.togglePlayPauseIcon();
                        self.toggleActiveTrack();
                        if (self.discs[discNumber-1].tracks.length === trackNumber
                            && self.discs.length === discNumber) {
                            return
                        }
                        var next = self.lastElement.parent().next();
                        if (next.length === 0) {
                            return
                        }
                        while (!next.hasClass('track')) {
                            next = next.next();
                            if (next.length === 0) {
                                return
                            }
                        }
                        next.children('.play').click()
                    }
                );
                self.lastElement = $(event.target);
                self.togglePlayPauseIcon();
                self.toggleActiveTrack();
            }
        };

        $('.play').click(self.playPause);

        self.player.preload(1, self.discs[0].tracks[0], function() {
            self.togglePlayPauseIcon();
            self.toggleActiveTrack();
            if (self.discs[0].tracks.length === 1
                && self.discs.length === 1) {
                return
            }
            var next = self.lastElement.parent().next();
            if (next.length === 0) {
                return
            }
            while (!next.hasClass('track')) {
                next = next.next();
                if (next.length === 0) {
                    return
                }
            }
            next.children('.play').click()
        });
    }

    window.app = new Album(album);
});