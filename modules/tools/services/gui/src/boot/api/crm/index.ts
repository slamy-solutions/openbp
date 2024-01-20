import { ClientAPI } from './client'
import { SettingsAPI } from './settings'
import { OneCAPI } from './onec'
import { DepartmentAPI } from './department'
import { ProjectAPI } from './project'
import { PerformerAPI } from './performer'
import { KanbanAPI } from './kanban'

export class CRMAPI {
    public readonly settings: SettingsAPI;
    public readonly client: ClientAPI;
    public readonly oneC: OneCAPI;
    public readonly department: DepartmentAPI;
    public readonly project: ProjectAPI;
    public readonly performer: PerformerAPI;
    public readonly kanban: KanbanAPI;

    constructor() {
        this.settings = new SettingsAPI()
        this.client = new ClientAPI()
        this.oneC = new OneCAPI()
        this.department = new DepartmentAPI()
        this.project = new ProjectAPI()
        this.performer = new PerformerAPI()
        this.kanban = new KanbanAPI()
    }
}