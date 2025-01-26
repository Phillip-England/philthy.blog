class TwMarkdown extends HTMLElement {
	constructor() {
	  super();
	}
  
	connectedCallback() {
	  const children = Array.from(this.children).map((child) =>
		child.cloneNode(true)
	  );
	  this.innerHTML = "";
	  children.forEach(this.styleElement);
	  children.forEach((child) => this.appendChild(child));
	}
  
	styleElement = (element) => {
	  const nodeName = element.nodeName.toLowerCase();
  
	  switch (nodeName) {
		case "pre":
		  element.classList.add(
			"custom-scroll",
			"p-4",
			"text-sm",
			"overflow-x-auto",
			"rounded",
			"mb-4",
		  );
		  break;
		case "h1":
		  element.classList.add(
			"font-bold",
			"text-3xl",
			"pb-4",
		  );
		  break;
		case "h2":
		  element.classList.add(
			"font-bold",
			"text-2xl",
			"pb-4",
			"pt-4",
			"border-t",
            "border-gray-200",
			"dark:border-gray-800",
		  );
		  break;
		case "h3":
		  element.classList.add(
			"font-bold",
			"text-xl",
			"mt-6",
			"mb-4",
		  );
		  break;
		case "p":
          let parent = element.parentElement;
          let nodeName = null
          if (parent != null) {
            nodeName = parent.nodeName.toLowerCase()            
          }
          if (nodeName == null) {
              element.classList.add(
                "text-sm",
                "leading-6",
                "mb-4",
              );
          } 
          if (nodeName == "blockquote") {
            element.classList.add(
                "text-sm",
                "leading-6",
              );
          }
		  break;
		case "ul":
		  element.classList.add(
			"pl-6",
			"mb-4",
			"list-disc",
		  );
		  break;
		case "ol":
		  element.classList.add(
			"pl-6",
			"mb-4",
			"list-decimal",
		  );
		  break;
		case "li":
		  element.classList.add(
			"mb-2",
			"text-sm",
		  );
		  break;
		case "blockquote":
		  element.classList.add(
			// "pl-1",
			// "border-l-1",
            "bg-gray-200",
            "dark:bg-dracula-background",
            "w-fit",
            "p-4",
            "rounded",
			// "border-gray-300",
			"italic",
			"text-gray-800",
            "dark:text-gray-200",
            "mb-4",
		  );
		  break;
		case "code":
		  if (element.parentElement.nodeName.toLowerCase() !== "pre") {
			element.classList.add(
			  "font-mono",
			  "px-1",
			  "rounded",
			  "text-sm",
			  "border",
			  "border-gray-200",
			  "dark:border-gray-800"
			);
		  }
		  break;
		case "hr":
		  element.classList.add(
			"border-t",
			"border-gray-300",
			"dark:border-gray-800",
			"my-4",
		  );
		  break;
		case "a":
		  element.classList.add(
			"text-blue-800",
			"underline",
            "visited:text-purple-500",
		  );
		  break;
		case "img":
		  element.classList.add(
			"max-w-full",
			"h-auto",
			"rounded",
			"my-4",
		  );
		  break;
	  }
  
	  // Recursively style child elements
	  Array.from(element.children).forEach(this.styleElement);
	};
  }
  
  class RandomBeads extends HTMLElement {
	connectedCallback() {
	  this.classList.add("flex", "flex-row", "gap-2");
	  const count = this.getAttribute("count");
	  const countInt = parseInt(count);
	  if (isNaN(countInt)) {
		console.error(
		  '<random-beads> requires an integer in the "count" attribute',
		);
		return;
	  }
	  this.beads = [];
	  let size = 4;
	  for (let i = 0; i < countInt; i++) {
		const bead = document.createElement("div");
		bead.classList.add("rounded-full", "transition-colors", "duration-1000");
		const initialColors = this.generateRandomColor();
		bead.style.height = `${size}px`;
		bead.style.width = `${size}px`;
		bead.style.backgroundColor =
		  `rgb(${initialColors.r}, ${initialColors.g}, ${initialColors.b})`;
  
		this.appendChild(bead);
		this.beads.push(bead);
		size += 1;
	  }
	  this.colorIntervalId = setInterval(() => this.transitionBeadColors(), 2000);
	}
	generateRandomColor() {
	  return {
		r: Math.floor(Math.random() * 256),
		g: Math.floor(Math.random() * 256),
		b: Math.floor(Math.random() * 256),
	  };
	}
	transitionBeadColors() {
	  this.beads.forEach((bead) => {
		const newColors = this.generateRandomColor();
		bead.style.backgroundColor =
		  `rgb(${newColors.r}, ${newColors.g}, ${newColors.b})`;
	  });
	}
	disconnectedCallback() {
	  if (this.colorIntervalId) {
		clearInterval(this.colorIntervalId);
	  }
	}
  }
  
  class TheBlinker extends HTMLElement {
	constructor() {
	  super();
	  this.attachShadow({ mode: "open" });
	}
	connectedCallback() {
	  const rate = parseInt(this.getAttribute("rate") || "1000");
	  const blinkElement = document.createElement("span");
	  blinkElement.textContent = this.textContent || "_";
	  const style = document.createElement("style");
	  style.textContent = `
		  @keyframes blink {
			0%, 100% { opacity: 1; }
			50% { opacity: 0; }
		  }
		  span {
			animation: blink ${rate}ms step-end infinite;
		  }
		`;
	  this.shadowRoot.appendChild(style);
	  this.shadowRoot.appendChild(blinkElement);
	}
  }
  
  class TitleLinks extends HTMLElement {
	constructor() {
	  super();
	}
  
	connectedCallback() {
	  const targetSelector = this.getAttribute("target");
	  const linkClass = this.getAttribute("link-class");
	  const linkWrapperClass = this.getAttribute("link-wrapper-class");
	  const linkClasses = linkClass.split(" ");
	  const linkWrapperClasses = linkWrapperClass.split(" ");
	  const offset = parseInt(this.getAttribute("offset"), 10) || 0;
  
	  const targetElement = document.querySelector(targetSelector);
	  if (!targetElement) {
		console.error(`Target element "${targetSelector}" not found.`);
		return;
	  }
  
	  const headings = targetElement.querySelectorAll("h1, h2, h3, h4, h5, h6");
  
	  headings.forEach((heading) => {
		if (heading.id) {
		  const linkItem = document.createElement("div");
		  linkWrapperClasses.forEach((linkWrapperClass) => {
			linkItem.classList.add(linkWrapperClass);
		  });
		  const link = document.createElement("a");
		  linkClasses.forEach((linkClass) => {
			link.classList.add(linkClass);
		  });
		  link.classList.add("title-link");
		  link.href = `#${heading.id}`;
		  link.textContent = heading.textContent;
		  linkItem.appendChild(link);
		  this.appendChild(linkItem);
		}
	  });
  
	  // Add styles
	  const style = document.createElement("style");
	  this.appendChild(style);
	  this.addEventListener("click", (e) => {
		if (e.target.tagName === "A") {
		  e.preventDefault();
		  const targetId = e.target.getAttribute("href").substring(1);
		  history.pushState({}, document.title, window.location.pathname+"#"+targetId)
		  const targetElement = document.getElementById(targetId);
		  if (targetElement) {
			const position = targetElement.getBoundingClientRect().top +
			  window.pageYOffset + offset;
			window.scrollTo({
			  top: position,
			  behavior: "smooth",
			});
		  }
		}
	  });
	}
  }
  
  class CustomScroll extends HTMLElement {
	constructor() {
	  super();
	}
	connectedCallback() {
	  this.innerHTML = `
			  <style>
				  .custom-scroll::-webkit-scrollbar {
					  width: 8px;
					  height: 8px;
				  }
				  .custom-scroll::-webkit-scrollbar-thumb {
					  background-color: #4B5563; /* Gray-600 */
					  border-radius: 4px;
				  }
				  .custom-scroll::-webkit-scrollbar-track {
					  background-color: #1F2937; /* Gray-800 */
				  }
				  /* Custom CSS to hide the scrollbar */
				  .scrollbar-hidden::-webkit-scrollbar {
					display: none;
				  }
  
				  .scrollbar-hidden {
					-ms-overflow-style: none;  /* For Internet Explorer and Edge */
					scrollbar-width: none;     /* For Firefox */
				  }
			  </style>
		  `;
	}
  }
  
  class HashTitleScroll extends HTMLElement {
	connectedCallback() {
	  let offset = parseInt(this.getAttribute('offset'), 10) || 0
	  let currentHref = window.location.href
	  let parts = currentHref.split('/')
	  let lastPart = parts[parts.length-1]
	  if (!lastPart.includes('#')) {
		return
	  }
	  let titleId = lastPart.split('#')[1]
	  let titleElm = document.getElementById(titleId)
	  if (!titleElm) {
		return
	  }
	  const position = titleElm.getBoundingClientRect().top + window.scrollY + offset;
	  window.scrollTo({
		top: position,
		behavior: "smooth"
	  })
	}
  }

  class BibleQuote extends HTMLElement {
    constructor() {
      super();
      this.title = this.getAttribute('title') || 'Verse';
      this.translation = this.getAttribute('translation') || 'Translation';
      this.verse = this.innerHTML.trim() || 'Verse text goes here.';
    }
  
    connectedCallback() {
      this.render();
    }
  
    render() {
      this.innerHTML = `
        <div class="bible-quote p-4 border border-gray-300 dark:border-dracula-background rounded mb-4 text-gray-600 dark:text-gray-400">
          <div class="bible-quote-header mb-4">
            <h2 class="text-lg">${this.title}</h2>
            <p class="text-xs italic">(${this.translation})</p>
          </div>
          <div class="bible-quote-body">
            <p class="text-sm">${this.verse}</p>
          </div>
        </div>
      `;
    }
  }
  
  
  window.addEventListener("DOMContentLoaded", () => {
	customElements.define("the-blinker", TheBlinker);
	customElements.define("tw-markdown", TwMarkdown);
	customElements.define("random-beads", RandomBeads);
	customElements.define("title-links", TitleLinks);
	customElements.define('hash-title-scroll', HashTitleScroll)
	customElements.define('custom-scroll', CustomScroll)
    customElements.define('bible-quote', BibleQuote);
  });