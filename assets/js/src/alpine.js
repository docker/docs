import Alpine from 'alpinejs'
import collapse from '@alpinejs/collapse'
 
window.Alpine = Alpine

Alpine.plugin(collapse)
Alpine.store("showSidebar", false)

Alpine.start()
