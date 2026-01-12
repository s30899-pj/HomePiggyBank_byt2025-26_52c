document.addEventListener('alpine:init', () => {
    Alpine.store('nav', {
        path: window.location.pathname
    })
})