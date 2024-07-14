var s; //snake
var foodNum = 10;
var foods = [];
var scl = 40; //pixel size
var score = 0;
var Hscore = 0;
var level = 5; //level
var boxW, boxH;
var ox, oy, nx, ny;

function setup() { 
        createCanvas(scl * floor(windowWidth / scl), scl * floor(windowHeight / scl));
        //console.log("windowWidth", windowWidth, "width", width);

        boxW = width - 2 * scl;
        boxH = height - 3 * scl;

        //console.log(boxW,"boxW");

        s = new Snake();
        strokeWeight(scl * 0.15);
        picklocation();
}

function mousePressed() {
        s.total ++;
}


function draw() {
        frameRate(level);
        background(200);

        drawGrid();
        drawFood();

        //draw snake
        s.update();
        s.show();
        s.death();
        strokeWeight(scl * 0.15);
        drawGame();
}


function pickFoodLocation() {
        let cols = floor(boxW / scl);
        let rows = floor(boxH / scl);
        let tmp  = createVector(1 + floor(random(cols)), 2 + floor(random(rows)));
        tmp.mult(scl);
        return tmp;
}

function picklocation() {       
        for (let i =0; i < foodNum ; i++ ) {  
                foods.push(pickFoodLocation());                              
        }
}

function keyPressed() {
        if (keyCode == UP_ARROW) {
                s.dir(0, -1);
        } else if (keyCode == DOWN_ARROW) {
                s.dir(0, 1);
        } else if (keyCode == LEFT_ARROW) {
                s.dir(-1, 0);
        } else if (keyCode == RIGHT_ARROW) {
                s.dir(1, 0);
        }
}

function drawGrid() {
        fill(190);
        stroke(200);
        for (let y = 2 * scl; y < height - scl; y += scl) {
                for (let x = scl; x < width - scl; x += scl) {
                        rect(x, y, scl, scl)
                }
        }
}

function drawGame() {
        // draw game box
        noFill();
        stroke(51)
        rect(scl, 2 * scl, boxW, boxH)

        //draw score
        fill(51);
        noStroke()
        textSize(0.7 * scl);
        text('Score: ' + score, scl, 0.8 * scl)
        text('High Score: ' + Hscore, scl, 1.55 * scl)
}


function drawFood() {        
        //fill(40, 20, 255, 150);        
        fill(40, 140, 40, 150);        
        for(let i = 0; i < foods.length; i ++ ) {
                let food = foods[i];
                if (s.eat(food)) {
                        level += ((level/2) / level);
                        score += floor(level);
                        if (score > Hscore) {
                                Hscore = score;
                        }
                        foods[i] = pickFoodLocation();                        
                } else {
                        ellipse(food.x + scl * 0.5, food.y + scl * 0.5, scl, scl);                        
                }
        }       
}



function touchStarted() {
        ox = mouseX;
        oy = mouseY;
}

function touchEnded() {
        nx = mouseX;
        ny = mouseY;

        if (abs((nx - ox) / (ny - oy)) > 1.5) {
                if (ox > nx) {
                        s.dir(-1, 0);
                } else {
                        s.dir(1, 0);
                }
        } else if (abs((ny - oy) / (nx - ox)) > 1.5) {
                if (oy < ny) {
                        s.dir(0, 1);
                } else {
                        s.dir(0, -1);
                }
        }
} 

function Snake() {
        this.x = 0;
        this.y = 0;
        this.xspeed = 1;
        this.yspeed = 0;
        this.total = 0;
        this.tail = [];

        this.death = function () {
                for (var i = 0; i < this.tail.length; i++) {
                        var pos = this.tail[i];
                        var d = dist(this.x, this.y, pos.x, pos.y);
                        if (d < 1) {
                                score = 0;
                                level = 5
                                this.total = 0;
                                this.tail = [];
                        }
                }
        }

        this.update = function () {
                for (var i = 0; i < this.tail.length - 1; i++) {
                        this.tail[i] = this.tail[i + 1];
                }
                this.tail[this.total - 1] = createVector(this.x, this.y);

                this.x = this.x + this.xspeed * scl;
                this.y = this.y + this.yspeed * scl;


                if (this.x > boxW) {
                        this.x -= boxW;
                } else if (this.x < scl) {
                        this.x += boxW;
                } else if (this.y < 2 * scl) {
                        this.y += boxH;
                } else if (this.y > boxH + scl) {
                        this.y -= boxH;
                }
                //this.x = constrain(this.x, scl, boxW);
                //this.y = constrain(this.y, 2 * scl, boxH);
        }

        this.show = function () {
                fill(255);
                fill(80);
                stroke(200);
                strokeWeight(scl * 0.15);

                for (var i = 0; i < this.total; i++) {
                        rect(this.tail[i].x, this.tail[i].y, scl, scl);
                }
                rect(this.x, this.y, scl, scl);
        }

        this.eat = function (pos) {
                var d = dist(this.x, this.y, pos.x, pos.y);
                if (d < 1) {
                        this.total++;
                        return true;
                } else {
                        return false;
                }
        }

        this.dir = function (x, y) {
                this.xspeed = x
                this.yspeed = y
        }
}
