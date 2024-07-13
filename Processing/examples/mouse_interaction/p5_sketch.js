let bubbles = [];
let bubblesNum = 20

function setup() {
        createCanvas(600, 400);
        /*
        for (let i = 0; i < bubblesNum; i++) {
                let x = random(width)
                let y = random(width)
                let r = random(10, 50)
                let b = new Bubble(x, y, r);
                bubbles.push(b);
        }
        */
}


function mouseDragged() {
        let r = random(10, 50)
        let b = new Bubble(mouseX, mouseY, r);
        bubbles.push(b);
}

function mousePressed() {
        for (let i = bubbles.length - 1; i >= 0; i--) {
                if (bubbles[i].contains(mouseX, mouseY)) {
                        bubbles.splice(i, 1)
                }
        }
}

function draw() {
        background(0);
        for (let i = 0; i < bubbles.length; i++) {
                if (bubbles[i].contains(mouseX, mouseY)) {
                        bubbles[i].changeColor(255)
                } else {
                        bubbles[i].changeColor(0)
                }
                bubbles[i].move();
                bubbles[i].show()
        }
        if (bubbles.length > 10) {
                bubbles.splice(0, 0)
        }
}

class Bubble {
        constructor(x, y, r) {
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
                this.x = this.x + random(-2, 2);
                this.y = this.y + random(-2, 2);
        }

        show() {
                stroke(255);
                strokeWeight(4);
                fill(this.brightness, 124)
                ellipse(this.x, this.y, this.r * 2);
        }

}
