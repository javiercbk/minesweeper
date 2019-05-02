export default {
  bind(el, { value }) {
    let originalOpacity;
    if (el.__vOriginalOpacity) {
      originalOpacity = el.__vOriginalOpacity;
    } else {
      originalOpacity = 1;
      el.__vOriginalOpacity = originalOpacity;
    }
    el.style.opacity = value ? 0 : originalOpacity;
  },

  update(el, { value, oldValue }) {
    if (!value === !oldValue) return;
    let originalOpacity = el.style.opacity;
    if (el.__vOriginalOpacity) {
      originalOpacity = el.__vOriginalOpacity;
    }
    el.style.opacity = value ? 0 : originalOpacity;
  },

  unbind(el, binding, vnode, oldVnode, isDestroy) {
    if (!isDestroy) {
      el.style.opacity = el.__vOriginalOpacity;
    }
  }
};
