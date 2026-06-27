export class Env {
  public static BirdsEyeApiEndpoint(): string {
    return import.meta.env.BIRDSEYE_VITE_BIRDS_EYE_API_ENDPOINT
      ? import.meta.env.BIRDSEYE_VITE_BIRDS_EYE_API_ENDPOINT
      : "https://birds-eye-api.ts-soda.net";
  }
}
