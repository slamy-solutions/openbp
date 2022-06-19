import { Testing } from './scheme'
import fs from 'fs'
import path from 'path'

export default async () => {
    const entries = await fs.promises.readdir(__dirname)
    const stats = await Promise.all(entries.map((entry) => fs.promises.stat(path.join(__dirname, entry))))
    const subfolders = entries.filter((_, index) => stats[index].isDirectory())
    for (const folder of subfolders) {
        try {
            const testingSuite = await require(path.join(__dirname, folder, "testing.ts")) as Testing
            await testingSuite.setup?.call(null)
        } catch (e) {
            console.warn("Error importing module testing suite. " + e)
        }
    }
}