export enum Side {
  Long = 0,
  Short = 1,
}

export type Room = {
  username: string;
  width: number;
  center: number;
};

export type GameState = {
  view: string;
  room: Room;
  market: string;
  pnl: number;
  place: number;
};
