import { h, signal } from "./framework.js";

const [getCount, setCount] = signal(0);

const Counter = () =>
  h.div(
    {
      style: {
        display: "flex",
        flexDirection: "column",
        gap: "0.5rem",
        alignItems: "center",
      },
    },
    [
      h.span({ style: { fontSize: "1.5rem", fontWeight: "bold" } }, [getCount]),
      h.div({ style: { display: "flex", gap: "0.5rem" } }, [
        h.button(
          {
            onClick: () => setCount(getCount() + 1),
            style: {
              padding: "0.5rem 1rem",
              background: "green",
              color: "white",
              border: "none",
              borderRadius: "4px",
              cursor: "pointer",
            },
          },
          ["Increment"],
        ),
        h.button(
          {
            onClick: () => setCount(getCount() - 1),
            style: {
              padding: "0.5rem 1rem",
              background: "red",
              color: "white",
              border: "none",
              borderRadius: "4px",
              cursor: "pointer",
            },
          },
          ["Decrement"],
        ),
      ]),
    ],
  );

document.getElementById("app").append(Counter());
