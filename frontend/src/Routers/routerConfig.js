
import loadable from 'utils/loadable'

// import LoadingComponent from 'components/LoadingComponent'
//过场组件默认采用通用的，若传入了loading，则采用传入的过场组件

const UserManagement = loadable(() => import('pages/UserManagement'))
const Glance = loadable(() => import('pages/Glance'))
const JXServingOnPremises = loadable(()=>import('pages/JXServingOnPremises'))
const DataSet = loadable(()=>import('pages/DataSet'))
const DatasetMange = loadable(() => import('pages/DatasetMange'))
const Label = loadable(()=>import('pages/Label'))
const Model = loadable(()=>import('pages/Model'))
const Analyst = loadable(()=>import('pages/Analyst'))
const Intervention = loadable(()=>import('pages/Intervention'))

const routerConfig = [
    {
        path:'/usermanagement',
        component:UserManagement
    },
    {
        path:'/glance',
        component:Glance,
        exact:true
    },
    {
        path:'/jxservingonpremises',
        component:JXServingOnPremises,
        exact:true
    },
    {
        path:'/dataset',
        component:DataSet,
    },
    {
        path:'/datasetmanage',
        component:DatasetMange
    },
    {
        path:'/label',
        component:Label
    },
    {
        path:'/model',
        component:Model,
    },
    {
        path:'/analyst',
        component:Analyst,
        exact:true
    },
    {
        path:'/intervention',
        component:Intervention,
    },
]
export default routerConfig