export interface LapSector {
  sectorNumber: number;
  sectorTime: number;
  isActive: boolean;
}

export interface Lap {
  lapNumber: number;
  lapTime: number;
  isValid: boolean;
  isActive: boolean;
  lapSectors: LapSector[];
}

export interface Session {
  id: string;
  startTime: number;
  sessionType: string;
  track: string;
  carModel: string;
  numberOfSectors: number;
  completedLaps: number;
  bestLapTime: number;
  previousLapTime: number;
  isActive: boolean;
  player: string;
  laps: Lap[];
}
