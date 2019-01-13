let state = new Map();
let states = new Map();

function clickItem(event) {
  let target = event.target;
  let style = window.getComputedStyle(target, null);

  if (style.opacity > 0.5) {
    let nStates = states.get(target).length;

    if (nStates > 0) {
      // cycle through progressive states
      if (state.get(target) == nStates - 1) {
        state.set(target, 0);
        target.style.opacity = 0.25;
      } else {
        state.set(target, state.get(target) + 1);
      }

      target.src = states.get(target)[state.get(target)];
    } else {
      // item is not progressive
      target.style.opacity = 0.25;
    }
  } else {
    target.style.opacity = 0.75;
  }
}

function init() {
  let items = document.getElementsByClassName("item");
  for (let i = 0; i < items.length; i++) {
    let item = items.item(i);
    item.onclick = clickItem;
    state.set(item, 0);
    states.set(item, []);
  }

  let sword = document.getElementById("sword");
  states.set(sword, ["img/sword1.gif", "img/sword2.gif", "img/sword3.gif"]);

  let flute = document.getElementById("flute");
  states.set(flute, ["img/flutestrange.gif", "img/flutericky.gif",
    "img/flutedimitri.gif", "img/flutemoosh.gif"]);

  let shield = document.getElementById("shield");
  states.set(shield, ["img/shield1.gif", "img/shield2.gif", "img/shield3.gif"]);
}
