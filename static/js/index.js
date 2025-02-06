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
        let nodeName = null;
        if (parent != null) {
          nodeName = parent.nodeName.toLowerCase();
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
            "dark:border-gray-800",
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
        element.setAttribute('target', '_blank')
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
        history.pushState(
          {},
          document.title,
          window.location.pathname + "#" + targetId,
        );
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
    let offset = parseInt(this.getAttribute("offset"), 10) || 0;
    let currentHref = window.location.href;
    let parts = currentHref.split("/");
    let lastPart = parts[parts.length - 1];
    if (!lastPart.includes("#")) {
      return;
    }
    let titleId = lastPart.split("#")[1];
    let titleElm = document.getElementById(titleId);
    if (!titleElm) {
      return;
    }
    const position = titleElm.getBoundingClientRect().top + window.scrollY +
      offset;
    window.scrollTo({
      top: position,
      behavior: "smooth",
    });
  }
}

class BibleQuote extends HTMLElement {
  constructor() {
    super();
    this.title = this.getAttribute("title") || "Verse";
    this.translation = this.getAttribute("translation") || "Translation";
    this.verse = this.innerHTML.trim() || "Verse text goes here.";
  }

  connectedCallback() {
    this.render();
  }

  render() {
    this.innerHTML = `
        <div class="bible-quote p-4 border border-gray-300 dark:border-dracula-background rounded mb-4 text-gray-800 dark:text-gray-400">
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
  customElements.define("hash-title-scroll", HashTitleScroll);
  customElements.define("custom-scroll", CustomScroll);
  customElements.define("bible-quote", BibleQuote);
});

//===============================
// classes
//===============================

class FullScreenToggle {
    constructor(onSelector, offSelector) {
        this.on = document.querySelector(onSelector)
        this.off = document.querySelector(offSelector)
        this.hook()
    }
    hook() {
        this.on.addEventListener('click', () => {
            this.toggleButtons()
            let success = this.enableFullScreen()
            if (success == false) {
                this.toggleButtons()
            }
        })
        this.off.addEventListener('click', () => {
            this.toggleButtons()
            let success = this.exitFullScreen()
            if (success == false) {
                this.toggleButtons()
            }
        })
    }
    toggleButtons() {
        this.on.classList.toggle('hidden')
        this.off.classList.toggle('hidden')
    }
    enableFullScreen() {
        // check for chrome
        if (document.documentElement.requestFullscreen) {
            document.documentElement.requestFullscreen();
            return true
        // check for safari
        } else if (document.documentElement.webkitRequestFullscreen) { // safari
            document.documentElement.webkitRequestFullscreen();
            return true
        // check for ie/edge
        } else if (document.documentElement.msRequestFullscreen) { // ie/edge
            document.documentElement.msRequestFullscreen();
            return true
        } else {
            console.warn('Fullscreen mode is not supported by this browser.');
            return false
        }
    }
    exitFullScreen() {
        // check for chrome
        if (document.exitFullscreen) {
            document.exitFullscreen();
            return true
        // check for safari
        } else if (document.webkitExitFullscreen) {
            document.webkitExitFullscreen();
            return true
        // check for ie/edge
        } else if (document.msExitFullscreen) {
            document.msExitFullscreen();
            return true
        } else {
            console.warn('Exiting full-screen mode is not supported by this browser.');
            return false
        }
    }
}

class AudioWhiteboard {
    constructor(recordSelector, stopSelector, whiteboardSelector) {
        this.record = document.querySelector(recordSelector);
        this.stop = document.querySelector(stopSelector);
        this.whiteboard = document.querySelector(whiteboardSelector);
        this.errorMode = false;
        this.shutdownMode = false;
        this.SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition;
        this.recognition = null;
        this.finalTranscript = '';
        this.retryCount = 0;
        this.maxRetries = 3;
        this.restartInterval = null;
        this.hookRecord();
        this.hookStop();
    }
    
    toggleButtons() {
        this.record.classList.toggle('hidden');
        this.stop.classList.toggle('hidden');
    }
    
    resetButtons() {
        this.record.classList.remove('hidden');
        this.stop.classList.add('hidden');
    }
    
    clearWhiteboard() {
        this.whiteboard.innerHTML = ``;
    }
    
    setError(msg) {
        this.whiteboard.innerHTML = `<p class='text-red-500'>${msg}</p>`;
        console.warn(msg);
    }
    
    startRecognition() {
        if (!this.SpeechRecognition) {
            this.errorMode = true;
            this.setError("Speech recognition not supported in this browser.");
            setTimeout(() => {
                this.clearWhiteboard();
                this.resetButtons();
                this.errorMode = false;
            }, 2000);
            return;
        }
        
        this.recognition = new this.SpeechRecognition();
        this.recognition.continuous = true;
        this.recognition.interimResults = true;
        this.recognition.lang = 'en-US';
        
        this.recognition.onresult = (event) => {
            if (this.shutdownMode) return;
            
            let interimTranscript = '';
            for (let i = event.resultIndex; i < event.results.length; i++) {
                const result = event.results[i];
                if (result.isFinal) {
                    this.finalTranscript += result[0].transcript.trim() + ' ';
                } else {
                    interimTranscript += result[0].transcript;
                }
            }
            
            const fullTranscript = (this.finalTranscript + interimTranscript)
                .replace(/\bfilthy\b/gi, 'philthy')
                .replace(/\bphilip\b/gi, 'phillip');
            
            const words = fullTranscript.trim().split(' ');
            const lastWord = words.pop();
            this.whiteboard.innerHTML = `${words.join(' ')} <span class="text-red-500">${lastWord}</span>`;
            
            window.scrollTo({ top: document.body.scrollHeight, behavior: 'smooth' });
        };
        
        this.recognition.onerror = (event) => {
            this.retryCount++
            this.errorMode = true;
            this.setError("Do you have a microphone plugged in or some way to capture audio?");
            setTimeout(() => {
                this.clearWhiteboard();
                this.resetButtons();
                this.errorMode = false;
            }, 2000);
        };
        
        this.recognition.onend = () => {
            if (this.retryCount >= this.maxRetries) {
                console.log('hit max retries, ending..')
                return
            }
            if (!this.shutdownMode) {
                console.log('Restarting recognition to prevent timeout.');
                this.startRecognition();
            }
        };
        
        this.recognition.start();
        console.log('Speech recognition started. Speak into the microphone.');
    }
    
    hookRecord() {
        this.record.addEventListener('click', () => {
            if (this.errorMode) {
                console.warn('Cannot record in error mode, wait a second');
                return;
            }
        
            this.retryCount = 0
            this.clearWhiteboard();
            this.finalTranscript = '';
            this.toggleButtons();
            this.shutdownMode = false;
            this.startRecognition();
            
            this.restartInterval = setInterval(() => {
                if (this.recognition) {
                    console.log('Manually restarting recognition to avoid timeout.');
                    this.recognition.stop();
                }
            }, 4 * 60 * 1000); // Restart every 4 minutes
        });
    }
    
    hookStop() {
        this.stop.addEventListener('click', () => {
            if (this.errorMode) {
                console.warn('Cannot stop recording in error mode, wait a second');
                return;
            }
            
            this.shutdownMode = true;
            clearInterval(this.restartInterval);
            
            if (this.recognition) {
                this.recognition.stop();
                this.recognition = null;
            }
            
            this.resetButtons();
            this.clearWhiteboard();
            setTimeout(() => {
                this.clearWhiteboard();
            }, 200);
            
            console.log('Speech recognition has stopped.');
        });
    }
}

class DarkModeToggler {
    constructor(sunSelector, moonSelector) {
        this.sun = document.querySelector(sunSelector);
        this.moon = document.querySelector(moonSelector);
        this.init();
    }
    init() {
        document.documentElement.classList.toggle('dark', 
            localStorage.theme === 'dark' || 
            (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)
        );
        window.addEventListener("DOMContentLoaded", () => {
            if (this.sun && this.moon) {
                this.sun.addEventListener("click", () => this.setLightMode());
                this.moon.addEventListener("click", () => this.setDarkMode());
            }
        });
    }
    setLightMode() {
        localStorage.theme = "light";
        document.documentElement.classList.remove("dark");
    }
    setDarkMode() {
        localStorage.theme = "dark";
        document.documentElement.classList.add("dark");
    }
}

//===============================
// path specific code
//===============================

const path = window.location.pathname.replace('.html', '')

// runs on all page loads
new DarkModeToggler('#sun', '#moon')

if (path == "/screenplay") {
    new FullScreenToggle('#expand', '#compress')
    new AudioWhiteboard('#record', '#stop', '#whiteboard')
}







