import React, { useContext, useState } from "react";
import { StateContext } from "./Game";

const Center = () => {
  const state = useContext(StateContext);
  const [center, setCenter] = useState<number>();

  if (state === undefined) {
    return <p>Loading...</p>;
  }

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (center === undefined) {
      return;
    }
    const formData = new FormData();
    formData.append("value", String(center));
    fetch("/center", {
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

  return (
    <div className="bg-white shadow-lg rounded-2xl p-8 w-96">
      <h2 className="text-2xl font-semibold text-center text-gray-800">
        Center Market
      </h2>

      <div className="mt-4 text-gray-700">
        <p className="text-lg">
          Market: <span className="font-semibold">{state.market}</span>
        </p>
        <p className="text-lg">
          Width: <span className="font-semibold">{state.room.width}</span>
        </p>
        <p className="text-lg">
          Current center:{" "}
          <span className="font-semibold">{state.room.center}</span>
        </p>
      </div>

      <form className="mt-6" onSubmit={handleSubmit}>
        <label className="block text-gray-600 font-medium mb-1">
          Center Market:
        </label>
        <input
          type="number"
          name="bid"
          placeholder="Enter center"
          value={center === undefined ? "" : String(center)}
          onChange={(e) => setCenter(Number(e.target.value))}
          min="0"
          autoComplete="off"
          required
          className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-400"
        />
        <button
          type="submit"
          className="w-full mt-4 bg-blue-500 text-white font-semibold py-2 rounded-lg hover:bg-blue-600 transition"
        >
          Submit
        </button>
      </form>
    </div>
  );
};

export default Center;
