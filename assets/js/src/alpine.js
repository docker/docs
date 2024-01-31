import Alpine from 'alpinejs'
import collapse from '@alpinejs/collapse'
import persist from '@alpinejs/persist'
 
window.Alpine = Alpine

Alpine.plugin(collapse)
Alpine.store("showSidebar", false)

Alpine.plugin(persist)

Alpine.start()
