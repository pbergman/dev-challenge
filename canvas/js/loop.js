var Loop = (function(f){
    var fps = ('undefined' === typeof f) ? 30 : f,
        frame = 0,
        interval,
        callable,
        initializeInterval = function(){
            interval = setInterval(function(){
                frame++;
                callable();
            }, 1000/fps);
        };

    return {
        getFramesPerSecond: function(){
            return fps;
        },
        setFramesPerSecond: function(f){
            fps = f;
        },
        getCurrentFrame: function(){
            return frame;
        },
        run: function(callback){

            if ('function' !== typeof callback) {
                throw new Error('The argument should be a callable function given: ' +  typeof callback);
            }

            callable = callback;

            initializeInterval();

        },
        reset: function(){
            clearInterval(interval);
            initializeInterval();
        },
        getPercentage: function(){
            return fps/1000*100;
        }
    }
});