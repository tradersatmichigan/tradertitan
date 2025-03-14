import React, { useContext, useState } from "react";
import { StateContext } from "./Game";
import { Side } from "./types";

const Trade = () => {
  const state = useContext(StateContext);
  const [side, setSide] = useState<Side>(Side.Long);

  if (state === undefined) {
    return <p>Loading...</p>;
  }

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData();
    formData.append("value", String(side));
    fetch("/trade", {
      method: "POST",
      body: formData,
    })
      .then(async (response) => {
        if (!response.ok) {
          throw new Error(await response.text());
        }
        return response.text();
      })
      .catch((error) => console.error("Fetch error:", error));
  };

  const buy_text = `BUY @ ${state.room.center + state.room.width / 2}`;
  const sell_text = `SELL @ ${state.room.center - state.room.width / 2}`;

  return (
    <>
      <p>Market: {state.market}</p>
      <form onSubmit={handleSubmit}>
        <span>Enter trade: </span>
        <label>
          <input
            type="radio"
            value="long"
            name="side"
            checked={side === Side.Long}
            onChange={(e) =>
              setSide(e.target.value === "long" ? Side.Long : Side.Short)
            }
          />
          {buy_text}
        </label>
        <label>
          <input
            type="radio"
            value="short"
            name="side"
            checked={side === Side.Short}
            onChange={(e) =>
              setSide(e.target.value === "short" ? Side.Short : Side.Long)
            }
          />
          {sell_text}
        </label>
        <button type="submit">Submit</button>
      </form>
    </>
  );
};

export default Trade;
