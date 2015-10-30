/**
 * Created by philip on 10/29/15.
 */

var CanvasObject = (function(element){

    if (element.tagName.toLowerCase() !== 'canvas') {
        throw new Error('Expeting canvas element got: "' + element.tagName.toLowerCase() + '"');
    }

    var cElement = element,
        cWidth = cElement.width,
        cHeight = cElement.height,
        rotation = 1,
        enabled = [
            'gray_circle',
            'z_letter',
            'blue_stripe'
        ],
        locked = false,
        drawings = {
            white_circle: function(ctx){
                ctx.save();
                ctx.fillStyle = "#ffffff";
                ctx.beginPath();
                ctx.moveTo(55.28,10.35);
                ctx.bezierCurveTo(68.06,9.89,81.01,14.87,89.96,24.05);
                ctx.bezierCurveTo(98.77,32.56,103.79,44.75,103.78,56.98);
                ctx.bezierCurveTo(103.8,69.24,98.76,81.46,89.93,89.98);
                ctx.bezierCurveTo(79.76,100.38,64.42,105.44,50.06,103.11);
                ctx.bezierCurveTo(37.21,101.24,25.37,93.61,18.2,82.79);
                ctx.bezierCurveTo(10.39,71.25,8.17,56.09,12.58,42.84);
                ctx.bezierCurveTo(18.22,24.51,36.1,10.89,55.28,10.35);
                ctx.closePath();
                ctx.fill();
                ctx.stroke();
                ctx.restore();
            },
            gray_circle: function(ctx){
                ctx.save();
                ctx.fillStyle = "#808084";
                ctx.beginPath();
                ctx.moveTo(50.21,2.27);
                ctx.bezierCurveTo(63.29,0.76,76.92,3.69,87.82,11.2);
                ctx.bezierCurveTo(101.21,20.02,110.29,35.01,111.94,50.96);
                ctx.bezierCurveTo(113.54,65,109.46,79.58,100.77,90.72);
                ctx.bezierCurveTo(90.12,104.64,72.54,113.01,54.99,112.02);
                ctx.bezierCurveTo(31.37,111.74,9.53,94.09,3.72,71.29);
                ctx.bezierCurveTo(0.34,58.26,1.63,43.92,7.88,31.93);
                ctx.bezierCurveTo(15.98,15.89,32.33,4.31,50.21,2.27);
                ctx.closePath();
                ctx.fill();
                ctx.stroke();
                ctx.restore();
                drawings['white_circle'](ctx);
            },
            z_letter: function(ctx) {
                ctx.save();
                ctx.fillStyle = "#808084";
                ctx.beginPath();
                ctx.moveTo(44.17,34.64);
                ctx.bezierCurveTo(39.98,34.06,40.23,27.58,44.09,26.65);
                ctx.bezierCurveTo(51.06,26.24,58.05,26.65,65.03,26.47);
                ctx.bezierCurveTo(67.37,26.61,70.24,25.97,72.04,27.89);
                ctx.bezierCurveTo(74.24,30.01,73.17,33.44,71.54,35.56);
                ctx.bezierCurveTo(65.21,44.39,58.83,53.18,52.49,62);
                ctx.bezierCurveTo(57.65,62,62.8,62.01,67.96,61.97);
                ctx.bezierCurveTo(69.89,62,72.12,61.89,73.6,63.38);
                ctx.bezierCurveTo(75.2,65.29,74.73,68.43,72.55,69.7);
                ctx.bezierCurveTo(70.2,70.82,67.49,70.44,64.96,70.53);
                ctx.bezierCurveTo(58.31,70.35,51.65,70.77,45.01,70.33);
                ctx.bezierCurveTo(42.28,70.41,39.98,67.72,40.54,65.04);
                ctx.bezierCurveTo(41.06,62.72,42.72,60.91,44.03,59);
                ctx.bezierCurveTo(49.82,50.94,55.6,42.87,61.35,34.78);
                ctx.bezierCurveTo(55.62,34.63,49.89,34.97,44.17,34.64);
                ctx.closePath();
                ctx.fill();
                ctx.stroke();
                ctx.restore();
            },
            blue_stripe: function(ctx){
                ctx.save();
                ctx.fillStyle = "#00adde";
                ctx.beginPath();
                ctx.moveTo(43.26,79.38);
                ctx.bezierCurveTo(52.49,79.09,61.76,79.27,71,79.29);
                ctx.bezierCurveTo(75.7,79.32,75.7,87.73,70.99,87.73);
                ctx.bezierCurveTo(61.99,87.78,52.99,87.77,43.99,87.73);
                ctx.bezierCurveTo(39.67,87.62,39.21,80.39,43.26,79.38);
                ctx.closePath();
                ctx.fill();
                ctx.stroke();
                ctx.restore();
            }
        },
        isEnabled = function(name){
                return enabled.indexOf(name) > -1;
        },
        start = function(ctx){
            ctx.save();
            ctx.beginPath();
            ctx.moveTo(0,0);
            ctx.lineTo(142.5,0);
            ctx.lineTo(142.5,142.5);
            ctx.lineTo(0,142.5);
            ctx.closePath();
            ctx.clip();
            ctx.scale(1.25,1.25);
            ctx.strokeStyle = 'rgba(0,0,0,0)';
            ctx.lineCap = 'butt';
            ctx.lineJoin = 'miter';
            ctx.miterLimit = 4;
        },
        end = function(ctx){
            ctx.restore();
        },
        clear = function(ctx){
            ctx.fillStyle = "#ffffff";
            ctx.clearRect(0, 0, cWidth, cHeight);
        },
        rotate = function(ctx){
            ctx.translate(cWidth/2, cHeight/2);
            ctx.rotate(rotation * Math.PI / 180);
            ctx.translate(-cWidth/2, -cHeight/2);
        };

    return {
        builder: function(){
            return {
                add: function(name, method){
                    drawings[name] = method;
                },
                remove: function(name){
                    delete drawings[name];

                }
            }
        },
        enable: function(name){
            if (!isEnabled(name)) {
                enabled.push(name);
            }
        },
        disable: function(name){
            var index = enabled.indexOf(name);
            if (index > -1) {
                enabled.splice(index, 1);
            }
        },
        render: function(loop){
            if (false === locked) {
                locked = true;
                ctx = cElement.getContext("2d");
                clear(ctx);
                rotate(ctx);
                start(ctx);
                enabled.forEach(function(name){
                    drawings[name](ctx);
                });
                end(ctx);
                locked = false;
            }
        },
        setRotation: function(v){
            rotation = v;
        }
    }
});
