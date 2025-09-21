/**
 * @template T
 * @param {T} initial
 * @returns {[get: (fn?: (v: T) => void) => T | (() => void), set: (v: T) => void]}
 */
export function signal(initial) {
  let value = initial;
  const subs = new Set();

  const get = (fn) => {
    if (fn) {
      subs.add(fn);
      fn(value); // fire immediately
      return () => subs.delete(fn); // unsubscribe
    }
    return value;
  };

  const set = (v) => {
    value = v;
    subs.forEach((fn) => fn(value));
  };

  return [get, set];
}

/**
 * @template T
 * @param {() => T} fn - Function producing the computed value
 * @param {Array<(fn: () => void) => void>} deps - Array of signals to subscribe to
 * @returns {() => T} - Getter for the computed value
 */
export function computed(fn, deps) {
  const [get, set] = signal(fn());

  const update = () => set(fn());
  deps.forEach((dep) => dep(update));

  return get;
}

/**
 * @param {() => void} fn - Side effect function
 * @param {Array<(fn: () => void) => void>} deps - Array of signals to subscribe to
 */
export function effect(fn, deps) {
  const run = () => fn();
  deps.forEach((dep) => dep(run));
  run();
}

/**
 * @param {string} tag
 * @returns {(propsOrChildren?: Record<string, any> | any[], children?: any[]) => HTMLElement}
 */
function el(tag) {
  return (arg, maybeChildren) => {
    /** @type {Record<string, any>} */
    let props = {};
    /** @type {any[]} */
    let children = [];

    if (Array.isArray(arg)) {
      children = arg;
    } else {
      props = arg || {};
      children = maybeChildren || [];
    }

    if (!Array.isArray(children)) {
      children = [children];
    }

    const node = document.createElement(tag);

    // props
    for (const [k, v] of Object.entries(props)) {
      if (k === "style" && typeof v === "object") {
        Object.assign(node.style, v);
      } else if (k.startsWith("on") && typeof v === "function") {
        node.addEventListener(k.slice(2).toLowerCase(), v);
      } else if (typeof v === "function") {
        v((val) => {
          node[k] = val;
        });
      } else {
        node[k] = v;
      }
    }

    // children
    for (let child of children.flat()) {
      if (child == null || child === false) continue;

      if (typeof child === "function") {
        let placeholder = document.createTextNode("");
        node.append(placeholder);

        child((val) => {
          const newNode = val?.nodeType
            ? val
            : document.createTextNode(val ?? "");
          placeholder.replaceWith(newNode);
          placeholder = newNode;
        });
      } else {
        node.append(child?.nodeType ? child : document.createTextNode(child));
      }
    }

    return node;
  };
}

/**
 * Hyperscript proxy — every property access gives you an element creator.
 * Example: h.div(props?, children?) → <div>
 * @type {Record<string, ReturnType<typeof el>>}
 */
export const h = new Proxy({}, { get: (_, tag) => el(tag) });
