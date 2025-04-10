class SlingShot {
  constructor(x, y, body) {
    const options = {
      pointA: {
        x: x,
        y: y,
      },
      bodyB: body,
      stiffness: 0.2,
      length: 40,
    };
    this.sling = Constraint.create(options);
    World.add(world, this.sling);
  }

  fly() {
    this.sling.bodyB = null;
  }

  show() {
    if (this.sling.bodyB) {
      push();
      stroke(255);
      strokeWeight(4);
      const posA = this.sling.pointA;
      const posB = this.sling.bodyB.position;
      line(posA.x, posA.y, posB.x, posB.y);
      pop();
    }
  }

  attach(body) {
    this.sling.bodyB = body;
  }
}
