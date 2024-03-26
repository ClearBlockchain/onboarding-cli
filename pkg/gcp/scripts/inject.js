const [, subscribeButton] = document.querySelectorAll(
  'button[data-prober="cloud-marketplace-request-product"]'
);

subscribeButton.style.margin = "0px";

const parent = subscribeButton.parentElement;

parent.style.padding = "10px";
parent.style.border = "2px solid rgba(35, 197, 82, 1)";
parent.style.borderRadius = "5px";
parent.style.position = "relative";
parent.style.transition = "border-color 0.5s";

let counter = 0;
const interval = setInterval(() => {
  if (counter % 2 === 0) {
    parent.style.borderColor = "rgba(35, 197, 82, 0.1)";
  } else {
    parent.style.borderColor = "rgba(35, 197, 82, 1)";
  }
  counter++;
  if (counter === 50) {
    clearInterval(interval);
  }
}, 500);

const textBubble = document.createElement("div");
textBubble.innerText = "Subscribe to ClearX OGI to complete the setup.";
textBubble.style.position = "absolute";
textBubble.style.top = "50%";
textBubble.style.left = "120%";
textBubble.style.transform = "translateY(-50%)";
textBubble.style.backgroundColor = "rgba(35, 197, 82, 1)";
textBubble.style.color = "#000";
textBubble.style.padding = "20px";
textBubble.style.borderRadius = "5px";
textBubble.style.fontSize = "20px";
textBubble.style.fontWeight = "400";
textBubble.style.zIndex = "1000";
textBubble.style.minWidth = "350px";
textBubble.style.textAlign = "center";
textBubble.style.boxShadow = "0px 0px 10px rgba(0, 0, 0, 0.1)";
textBubble.style.height = "auto";
textBubble.style.letterSpacing = "1px";
textBubble.style.lineHeight = "1.5";
parent.appendChild(textBubble);
