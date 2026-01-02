declare module '../bindings/github.com/XgzK/intellijapp/internal/service/configservice' {
  import { AboutInfo } from './models'

  export interface AssetInfo {
    name: string
    downloadUrl: string
    size: number
  }

  export interface ReleaseInfo {
    version: string
    publishedAt: string
    htmlUrl: string
    body: string
    assets: AssetInfo[]
  }

  export interface UpdateCheckResult {
    hasUpdate: boolean
    release: ReleaseInfo | null
  }

  export function SubmitPaths(projectPath: string, configPath: string): Promise<string>
  export function ClearConfig(projectPath: string): Promise<string>
  export function PathExists(path: string): Promise<boolean>
  export function GetAboutInfo(): Promise<AboutInfo>
  export function CheckForUpdates(): Promise<UpdateCheckResult>
  export function GetAccessibleGitHubMirror(): Promise<string>
  export function ConvertToAccessibleURL(originalURL: string): Promise<string>
}
