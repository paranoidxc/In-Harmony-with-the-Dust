let bubbles = [];
let bubblesNum = 20

let unicorn;

function setup() {
        createCanvas(800, 600);
        for (let i = 0; i < bubblesNum; i++) {
                bubbles[i] = new Bubble(random(width), random(height))
        }

        //unicorn = new Bubble(400, 200, 10)
}


function draw() {
        background(0);
        /*
        unicorn.x = mouseX
        unicorn.y = mouseY
        unicorn.show()
        unicorn.move()
        */

        for (let b of bubbles) {
                b.show()
                b.move()
                let overlapping = false
                for (let other of bubbles) {
                        if (b !== other && b.intersects(other)) {
                                overlapping = true
                                break;
                        }
                }
                if (overlapping) {
                        b.changeColor(255)
                } else {
                        b.changeColor(0)
                }
        }
}

class Bubble {
        constructor(x, y, r = 40) {
                this.x = x;
                this.y = y;
                this.r = r;
                this.brightness = 0;
        }

        changeColor(brightness) {
                this.brightness = brightness
        }

        contains(px, py) {
                let d = dist(px, py, this.x, this.y);
                return d < this.r;
        }

        clicked(px, py) {
                let d = dist(px, py, this.x, this.y);
                if (d < this.r) {
                        this.brightness = 255;
                        console.log("CLICKED");
                }
        }

        move() {
                this.x = this.x + random(-10, 10);
                this.y = this.y + random(-10, 10);
        }

        show() {
                stroke(255);
                strokeWeight(4);
                fill(this.brightness, 124)
                ellipse(this.x, this.y, this.r * 2);
        }

        intersects(other) {
                let d = dist(this.x, this.y, other.x, other.y);
                return d < (this.r + other.r)
        }
}
