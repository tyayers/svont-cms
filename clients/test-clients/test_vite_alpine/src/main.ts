import './style.css'
import typescriptLogo from './typescript.svg'
import { setupCounter } from './counter'
import Alpine from 'alpinejs'
 
window.Alpine = Alpine
Alpine.start()

document.querySelector<HTMLDivElement>('#app')!.innerHTML = `
  <div>
    <a href="https://vitejs.dev" target="_blank">
      <img src="/vite.svg" class="logo" alt="Vite logo" />
    </a>
    <a href="https://www.typescriptlang.org/" target="_blank">
      <img src="${typescriptLogo}" class="logo vanilla" alt="TypeScript logo" />
    </a>
    <h1>Vite + TypeScript</h1>
    <div class="card">
      <button id="counter" type="button"></button>
    </div>
    <p class="read-the-docs">
      Click on the Vite and TypeScript logos to learn more
    </p>

    <div x-data="{ open: false }">
        <button @click="open = ! open">Expand</button>
    
        <div x-show="open">
            Content...
        </div>
    </div>
  </div>
`

setupCounter(document.querySelector<HTMLButtonElement>('#counter')!)
