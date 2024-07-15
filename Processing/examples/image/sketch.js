let bubbles = [];
let bubblesNum = 20

let flower;
let kittens = [];

function preload() {
        flower = loadImage("./balloon1.png");
        for (let i = 0; i < 3; i++) {
                kittens[i] = loadImage("./balloon" + (i + 1) + ".png")
        }
}

function setup() {
        createCanvas(800, 600);
        preload()
        for (let i = 0; i < bubblesNum; i++) {
                let x = random(width)
                let y = random(width)
                let r = 50
                let idx = i % 3
                let b = new Bubble(x, y, r, kittens[idx]);
                bubbles.push(b);
        }
}


function mouseDragged() {
        /*
        let r = random(10, 50)
        let b = new Bubble(mouseX, mouseY, r, kittens[int(random(0,3))]);
        bubbles.push(b);
*/
}

function mousePressed() {
        for (let i = bubbles.length - 1; i >= 0; i--) {
                bubbles[i].clicked(mouseX, mouseY)
                /*
                if (bubbles[i].contains(mouseX, mouseY)) {
                        //bubbles.splice(i, 1)
                }
                */
        }
}

function draw() {
        background(0);
        for (let i = 0; i < bubbles.length; i++) {
                bubbles[i].move();
                bubbles[i].show()
        }
        /*
                if (bubbles[i].contains(mouseX, mouseY)) {
                        bubbles[i].changeColor(255)
                } else {
                        bubbles[i].changeColor(0)
                }
                        */
}

class Bubble {
        constructor(x, y, r, img) {
                this.x = x;
                this.y = y;
                this.r = r;
                this.kitten = img
        }

        changeColor(brightness) {
        }

        contains(px, py) {
                let d = dist(px, py, this.x, this.y);
                return d < this.r;
        }

        clicked(px, py) {
                if (px > this.x && px < this.x + this.r && py > this.y && py < this.y + this.r) {
                        this.kitten = random(kittens)
                }
                /*
                let d = dist(px, py, this.x, this.y);
                if (d < this.r) {
                        console.log("CLICKED");
                }
                */
        }

        move() {
                this.x = this.x + random(-2, 2);
                this.y = this.y + random(-2, 2);
        }

        show() {
                image(this.kitten, this.x, this.y, this.r, this.r)
        }

}
