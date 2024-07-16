// The Nature of Code
// Daniel Shiffman
// http://natureofcode.com

let particleA;
let particleB;

function setup() {
  createCanvas(640, 360);
  particleA = new Particle(320, 60);
  particleB = new Particle(320, 300);
}

function draw() {
  background(0);

  particleA.collide(particleB);

  particleA.update();
  particleB.update();

  particleA.edges();
  particleB.edges();

  particleA.show();
  particleB.show();

  // let speedA = particleA.velocity.mag();
  // let speedB = particleB.velocity.mag();
  // let kinA = 0.5 * particleA.mass * speedA * speedA;
  // let kinB = 0.5 * particleB.mass * speedB * speedB;
  // console.log(kinA + kinB);
}
