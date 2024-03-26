import Alpine from 'alpinejs'
import collapse from '@alpinejs/collapse'
import focus from '@alpinejs/focus'
 
window.Alpine = Alpine

Alpine.plugin(collapse)
Alpine.plugin(focus)
Alpine.store("showSidebar", false)

Alpine.start()
