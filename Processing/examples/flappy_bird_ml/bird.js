class Bird {
  constructor(brain) {
    if (brain) {
      this.brain = brain;
    } else {
      this.brain = ml5.neuralNetwork({
        inputs: 4,
        outputs: ["flap", "no flap"],
        task: "classification",

        // change to "neuroEvolution" for next ml5.js release
        neuroEvolution: true,
      });
    }

    this.x = 50;
    this.y = 120;
    this.w = 16;
    this.h = 16;

    this.r = sqrt((this.w / 2) * (this.w / 2) + (this.h / 2) * (this.h / 2));
    console.log("bird r:", this.r);

    this.velocity = 0;
    this.gravity = 0.5;
    this.flapForce = -10;

    this.angle = 0;
    this.flapAngle = 0;

    this.fitness = 0;
    this.alive = true;
  }

  think(pipes) {
    let nextPipe = null;
    for (let pipe of pipes) {
      if (pipe.x + pipe.w > this.x - 8) {
        nextPipe = pipe;
        break;
      }
    }

    let inputs = [
      (this.y - this.r) / height,
      this.velocity / height,
      nextPipe.top / height,
      (nextPipe.x - this.x - this.r) / width,
    ];

    let results = this.brain.classifySync(inputs);
    if (results[0].label == "flap") {
      this.flap();
    }
  }

  flap() {
    this.velocity += this.flapForce;
    this.flapAngle = 10;
  }

  update() {
    this.velocity += this.gravity;
    this.y += this.velocity;
    this.velocity *= 0.95;

    if (this.flapAngle > 0) {
      this.angle += 20;
      this.flapAngle -= 1;
    } else {
      this.angle += 2;
    }

    if (this.y > height || this.y < 0) {
      this.alive = false;
    }

    this.fitness++;
  }

  show() {
    push();
    angleMode(DEGREES);
    strokeWeight(2);
    stroke(0);
    fill(127, 200);
    // works
    //rect(this.x - this.w / 2, this.y - this.h / 2, this.w, this.h);

    translate(this.x - this.w / 2, this.y - this.h / 2);
    rotate(this.angle);
    rect(-this.w / 2, -this.h / 2, this.w, this.h);

    stroke("red");
    circle(0, 0, this.r * 2);

    pop();
  }
}
