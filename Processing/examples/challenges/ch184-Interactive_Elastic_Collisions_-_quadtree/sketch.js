// The Nature of Code
// Daniel Shiffman
// http://natureofcode.com
//
// Edited by:
// Radek Titěra
//
// Click on a particle to select it
// if mouse is not over any particle, new particle is created
// Drag the selected particle to move it
// Release the selected particle to fling it
//
// Under the canvas are some simulation and drawing settings
//
// Have fun ;)


let particles = [];
let qtree;

let draggedParticle = null;
let fling = {};

let stoppedTime = false;

let drawSettings = {
  useFill: false,
  drawVectors: false,
  redrawDragged: true,
};

// Making sure pairs of particles are not checked twice
let checkedPairs = new Set();
let palette = [
  [11, 106, 136],
  [45, 197, 244],
  [112, 50, 126],
  [146, 83, 161],
  [164, 41, 99],
  [236, 1, 90],
  [240, 99, 164],
  [241, 97, 100],
  [248, 158, 79],
  [252, 238, 33],
];

const minMass = 4;
const maxMass = 260;

let bgAlphaSlider, bgAlphaP;

let newParticleMassSlider, newParticleMassP;

function setup() {
  createCanvas(700, 500);

  fling = {
    start: createVector(),
    end: createVector(),
  };

  for (let i = 0; i < 50; i++) {
    let x = random(width);
    let y = random(height);
    let mass = random(minMass, maxMass);
    particles.push(new Particle(x, y, mass, i));
  }

  constructQuadTree();

  const stoppedTimeCheckbox = createCheckbox("Stop time?", stoppedTime);
  stoppedTimeCheckbox.mousePressed(function () {
    stoppedTime = !stoppedTimeCheckbox.checked();
  });
  
  newParticleMassP = createP("New particle mass:");
  newParticleMassSlider = createSlider(
    minMass,
    maxMass,
    (minMass + maxMass) / 2,
    0.1
  );
  newParticleMassP.html(`New particle mass: ${newParticleMassSlider.value()}`);
  
  createElement('hr');

  const useFillCheckbox = createCheckbox(
    "Fill particles?",
    drawSettings.useFill
  );
  useFillCheckbox.mousePressed(function () {
    drawSettings.useFill = !useFillCheckbox.checked();
  });

  const redrawDraggedCheckbox = createCheckbox(
    "Redraw dragged particle?",
    drawSettings.redrawDragged
  );
  redrawDraggedCheckbox.mousePressed(function () {
    drawSettings.redrawDragged = !redrawDraggedCheckbox.checked();
  });

  const drawVectorsCheckbox = createCheckbox(
    "Draw velocity vectors?",
    drawSettings.drawVectors
  );
  drawVectorsCheckbox.mousePressed(function () {
    drawSettings.drawVectors = !drawVectorsCheckbox.checked();
  });

  bgAlphaP = createP("Background alpha:");
  bgAlphaSlider = createSlider(0, 255, 32, 1);
  bgAlphaP.html(`Background alpha: ${bgAlphaSlider.value()}`);

}

function resolveCollisions() {
  // Create a quadtree
  checkedPairs.clear();

  for (let i = 0; i < particles.length; i++) {
    let particleA = particles[i];
    let range = new Circle(
      particleA.position.x,
      particleA.position.y,
      particleA.r * 2
    );

    // Check only nearby particles based on quadtree
    let points = qtree.query(range);
    for (let point of points) {
      let particleB = point.userData;
      if (particleB !== particleA) {
        let idA = particleA.id;
        let idB = particleB.id;
        let pair = idA < idB ? `${idA},${idB}` : `${idB},${idA}`;
        if (!checkedPairs.has(pair)) {
          particleA.collide(particleB);
          checkedPairs.add(pair);
        }
      }
    }
  }
}

function constructQuadTree() {
  let boundary = new Rectangle(width / 2, height / 2, width, height);
  qtree = new QuadTree(boundary, 4);

  // Insert all particles
  for (let particle of particles) {
    let point = new Point(particle.position.x, particle.position.y, particle);
    qtree.insert(point);
  }
}

function draw() {
  bgAlphaP.html(`Background alpha: ${bgAlphaSlider.value()}`);
  newParticleMassP.html(`New particle mass: ${newParticleMassSlider.value()}`);

  background(0, bgAlphaSlider.value());

  constructQuadTree();

  if (!stoppedTime) resolveCollisions();

  if (draggedParticle !== null) updateDraggedParticle();

  for (let particle of particles) {
    if (!stoppedTime) particle.update();
    particle.edges();
    particle.show();
  }

  if (frameCount % 120 == 0) {
    console.log(frameRate());
  }

  if (drawSettings.redrawDragged && draggedParticle !== null) {
    background(0, 150);
    draggedParticle.show();
  }
}

function trySetDraggedParticle() {
  // Find all points that might intersect the mouse
  const maxR = sqrt(maxMass) * 2;
  const range = new Circle(mouseX, mouseY, maxR);
  const points = qtree.query(range);

  let closest = null;
  let closestDistance = Number.MAX_SAFE_INTEGER;
  for (const point of points) {
    const particle = point.userData;
    const distance = dist(
      mouseX,
      mouseY,
      particle.position.x,
      particle.position.y
    );
    if (distance < particle.r && distance < closestDistance) {
      closest = particle;
      closestDistance = distance;
    }
  }
  draggedParticle = closest;

  if (draggedParticle !== null) {
    draggedParticle.velocity.mult(0);
    draggedParticle.acceleration.mult(0);
  }

  fling.start.set(mouseX, mouseY);
  fling.end.set(mouseX, mouseY);
}

function updateDraggedParticle() {
  fling.start = fling.end.copy();
  fling.end.set(mouseX, mouseY);

  const dragForce = p5.Vector.sub(fling.end, fling.start);
  dragForce.mult(0.5);

  draggedParticle.position.set(mouseX, mouseY);
  draggedParticle.velocity = dragForce;
}

function flingDraggedParticle() {
  const flingForce = p5.Vector.sub(fling.end, fling.start);
  flingForce.mult(80);

  draggedParticle.applyForce(flingForce);
}

function mousePressed() {
  trySetDraggedParticle();

  // If not over any particles and over canvas, create a new particle
  const mouseOverCanvas =
    mouseX >= 0 && mouseX < width && mouseY >= 0 && mouseY < height;
  if (mouseOverCanvas && draggedParticle === null) {
    const newParticle = new Particle(
      mouseX,
      mouseY,
      newParticleMassSlider.value(),
      particles.length
    );
    particles.push(newParticle);

    draggedParticle = newParticle;
  }
}

function mouseReleased() {
  if (draggedParticle !== null) {
    flingDraggedParticle();
    draggedParticle = null;
  }
}
