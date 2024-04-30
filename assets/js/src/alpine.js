import Alpine from 'alpinejs'
import collapse from '@alpinejs/collapse'
import persist from '@alpinejs/persist'
import focus from '@alpinejs/focus'
 
window.Alpine = Alpine

Alpine.plugin(collapse)
Alpine.plugin(persist)
Alpine.plugin(focus)
Alpine.store("showSidebar", false)

Alpine.start()
