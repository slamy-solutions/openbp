import { APIModuleBase } from '../model';
import { ListAPI } from './list'

export class NamespaceAPI extends APIModuleBase {
    public readonly list: ListAPI;

    constructor() {
        super();
        this.list = new ListAPI()
    }
}