
import React from 'react';
import {
    G2,
    Chart,
    Geom,
    Axis,
    Tooltip,
    Coord,
    Label,
    Legend,
    View,
    Guide,
    Shape,
    Facet,
    Util,
} from 'bizcharts';

const colors = ['rgba(33,144,233,1)','rgba(196,96,255,1)','rgba(227,23,13,1)','rgba(255,128,0,1)','rgba(255,215,0,1)','rgba(127,255,0,1)','rgba(188,143,143,1)','rgba(202,235,216,1)','rgba(255,245,238,1)','rgba(237,145,33,1)','rgba(199,97,20,1)','rgba(244,164,95,1)','rgba(210,180,140,1)','rgba(218,112,214,1)','rgba(0,199,140,1)','rgba(107,142,35,1)','rgba(64,224,205,1)','rgba(0,255,255,1)','rgba(3,168,158,1)','rgba(255,192,203,1)']
const thedata = [{ "id": "5d831810404ade3d94b8d124", "value": 0.4835672091497396, "type": "data" }, { "id": "5d831810404ade3d94b8d127", "value": 0.5423692647025724, "type": "data" }, { "id": "5d831810404ade3d94b8d12a", "value": 0.6666791747003366, "type": "data" }, { "id": "5d831810404ade3d94b8d133", "value": 0.36840747994821565, "type": "data" }, { "id": "5d831810404ade3d94b8d136", "value": 0.6650391460726056, "type": "data" }, { "id": "5d831810404ade3d94b8d142", "value": 0.37430259825121426, "type": "data" }, { "id": "5d831810404ade3d94b8d148", "value": 0.6874796996049942, "type": "data" }, { "id": "5d831810404ade3d94b8d14b", "value": 0.49814453504331657, "type": "data" }, { "id": "5d831810404ade3d94b8d154", "value": 0.6883747015794269, "type": "data" }, { "id": "5d831810404ade3d94b8d15a", "value": 0.5139333747442493, "type": "data" }, { "id": "5d831810404ade3d94b8d15d", "value": 0.6466337833688549, "type": "data" }, { "id": "5d831810404ade3d94b8d160", "value": 0.5081691225699947, "type": "data" }, { "id": "5d831810404ade3d94b8d16f", "value": 0.6438906377477148, "type": "data" }, { "id": "5d831810404ade3d94b8d172", "value": 0.6975273043646724, "type": "data" }, { "id": "5d831810404ade3d94b8d175", "value": 0.6376346911206766, "type": "data" }, { "id": "5d831810404ade3d94b8d178", "value": 0.32300643208654, "type": "data" }, { "id": "5d831810404ade3d94b8d18a", "value": 0.5672558590585323, "type": "data" }, { "id": "5d831810404ade3d94b8d18d", "value": 0.4813782226363862, "type": "data" }, { "id": "5d831810404ade3d94b8d190", "value": 0.6955457973917809, "type": "data" }, { "id": "5d831810404ade3d94b8d193", "value": 0.39772652661584473, "type": "data" }, { "id": "5d831810404ade3d94b8d1a8", "value": 0.4203697344226869, "type": "data" }, { "id": "5d831810404ade3d94b8d1b1", "value": 0.5360943536516244, "type": "data" }, { "id": "5d831810404ade3d94b8d1b4", "value": 0.49046782257076926, "type": "data" }, { "id": "5d831810404ade3d94b8d1bd", "value": 0.6656028543313761, "type": "data" }, { "id": "5d831810404ade3d94b8d1cc", "value": 0.37009466984249195, "type": "data" }, { "id": "5d831810404ade3d94b8d1d2", "value": 0.6894343824847702, "type": "data" }, { "id": "5d831810404ade3d94b8d1d5", "value": 0.4003191063952003, "type": "data" }, { "id": "5d831810404ade3d94b8d1db", "value": 0.39225360146493826, "type": "data" }, { "id": "5d831810404ade3d94b8d1ea", "value": 0.3762247707282169, "type": "data" }, { "id": "5d831811404ade3d94b8d1ed", "value": 0.395547593668209, "type": "data" }, { "id": "5d831811404ade3d94b8d1f0", "value": 0.7250092461673668, "type": "data" }, { "id": "5d831811404ade3d94b8d1f3", "value": 0.6023471074549585, "type": "data" }, { "id": "5d831811404ade3d94b8d1ff", "value": 0.6154860055850991, "type": "data" }, { "id": "5d831811404ade3d94b8d202", "value": 0.3313783169835933, "type": "data" }, { "id": "5d831811404ade3d94b8d208", "value": 0.34826051857950235, "type": "data" }, { "id": "5d831811404ade3d94b8d217", "value": 0.5560349505702318, "type": "data" }, { "id": "5d831811404ade3d94b8d21d", "value": 0.38624453511794343, "type": "data" }, { "id": "5d831811404ade3d94b8d223", "value": 0.47511477361013255, "type": "data" }, { "id": "5d831811404ade3d94b8d226", "value": 0.5244158120768063, "type": "data" }, { "id": "5d831811404ade3d94b8d232", "value": 0.39732794755699197, "type": "data" }, { "id": "5d831811404ade3d94b8d241", "value": 0.39346706832906897, "type": "data" }, { "id": "5d831812404ade3d94b8d250", "value": 0.6124107734398495, "type": "data" }, { "id": "5d831812404ade3d94b8d256", "value": 0.35416580149119015, "type": "data" }, { "id": "5d831812404ade3d94b8d259", "value": 0.3769016724968053, "type": "data" }, { "id": "5d831812404ade3d94b8d25f", "value": 0.5632848053097376, "type": "data" }, { "id": "5d831812404ade3d94b8d262", "value": 0.6970098888738091, "type": "data" }, { "id": "5d831812404ade3d94b8d265", "value": 0.7201917642221619, "type": "data" }, { "id": "5d831812404ade3d94b8d268", "value": 0.5634909336609634, "type": "data" }, { "id": "5d831812404ade3d94b8d26b", "value": 0.6887614818747204, "type": "data" }, { "id": "5d831812404ade3d94b8d277", "value": 0.6350118990127102, "type": "data" }, { "id": "5d831812404ade3d94b8d27d", "value": 0.34456593099808064, "type": "data" }, { "id": "5d831812404ade3d94b8d280", "value": 0.3616712740697504, "type": "data" }, { "id": "5d831812404ade3d94b8d283", "value": 0.3351169277184135, "type": "data" }, { "id": "5d831812404ade3d94b8d286", "value": 0.4192981145237795, "type": "data" }, { "id": "5d831812404ade3d94b8d289", "value": 0.6489491472977281, "type": "data" }, { "id": "5d831812404ade3d94b8d28c", "value": 0.3251657974623219, "type": "data" }, { "id": "5d831812404ade3d94b8d28f", "value": 0.6580768970033805, "type": "data" }, { "id": "5d831812404ade3d94b8d296", "value": 0.3734056629347505, "type": "data" }, { "id": "5d831812404ade3d94b8d299", "value": 0.7052090062153881, "type": "data" }, { "id": "5d831812404ade3d94b8d29c", "value": 0.6441214451693957, "type": "data" }, { "id": "5d831812404ade3d94b8d2a2", "value": 0.5331437867399289, "type": "data" }, { "id": "5d831812404ade3d94b8d2a8", "value": 0.479815516402192, "type": "data" }, { "id": "5d831812404ade3d94b8d2ab", "value": 0.7203496113206816, "type": "data" }, { "id": "5d831812404ade3d94b8d2b1", "value": 0.7282753159532952, "type": "data" }, { "id": "5d831812404ade3d94b8d2b4", "value": 0.4143838395329984, "type": "data" }, { "id": "5d831812404ade3d94b8d2b7", "value": 0.5375592025841816, "type": "data" }, { "id": "5d831812404ade3d94b8d2ba", "value": 0.3419043422417216, "type": "data" }, { "id": "5d831812404ade3d94b8d2c3", "value": 0.5225319554198097, "type": "data" }, { "id": "5d831812404ade3d94b8d2cc", "value": 0.5163647010254215, "type": "data" }, { "id": "5d831812404ade3d94b8d2cf", "value": 0.528197818494487, "type": "data" }, { "id": "5d831812404ade3d94b8d2d5", "value": 0.6273545261511063, "type": "data" }, { "id": "5d831812404ade3d94b8d2e1", "value": 0.389557299618589, "type": "data" }, { "id": "5d831812404ade3d94b8d2e7", "value": 0.47680907512533455, "type": "data" }, { "id": "5d831812404ade3d94b8d2ea", "value": 0.44309195925436196, "type": "data" }, { "id": "5d831812404ade3d94b8d2f0", "value": 0.7120291609541165, "type": "data" }, { "id": "5d831812404ade3d94b8d2f6", "value": 0.6823309477489261, "type": "data" }, { "id": "5d831812404ade3d94b8d2ff", "value": 0.3606892559931537, "type": "data" }, { "id": "5d831812404ade3d94b8d302", "value": 0.6779266995582461, "type": "data" }, { "id": "5d831812404ade3d94b8d305", "value": 0.4145418245658177, "type": "data" }, { "id": "5d831812404ade3d94b8d308", "value": 0.39472316890149683, "type": "data" }, { "id": "5d831812404ade3d94b8d30b", "value": 0.6619490518634128, "type": "data" }, { "id": "5d831812404ade3d94b8d317", "value": 0.4940487120590584, "type": "data" }, { "id": "5d831812404ade3d94b8d320", "value": 0.6396142935912091, "type": "data" }, { "id": "5d831812404ade3d94b8d323", "value": 0.435651650562057, "type": "data" }, { "id": "5d831812404ade3d94b8d326", "value": 0.5175180701113756, "type": "data" }, { "id": "5d831812404ade3d94b8d329", "value": 0.4938107119365368, "type": "data" }, { "id": "5d831812404ade3d94b8d32f", "value": 0.7084076308825001, "type": "data" }, { "id": "5d831812404ade3d94b8d332", "value": 0.46626589583477585, "type": "data" }, { "id": "5d831812404ade3d94b8d335", "value": 0.5019242675605575, "type": "data" }, { "id": "5d831812404ade3d94b8d338", "value": 0.6695090514522187, "type": "data" }, { "id": "5d831812404ade3d94b8d33b", "value": 0.6051932265632953, "type": "data" }, { "id": "5d831812404ade3d94b8d347", "value": 0.36247286646416066, "type": "data" }, { "id": "5d831812404ade3d94b8d34d", "value": 0.6219184440981742, "type": "data" }, { "id": "5d831812404ade3d94b8d350", "value": 0.4479792806254922, "type": "data" }, { "id": "5d831812404ade3d94b8d362", "value": 0.3767232558208581, "type": "data" }, { "id": "5d831812404ade3d94b8d371", "value": 0.5666688059353515, "type": "data" }, { "id": "5d831813404ade3d94b8d374", "value": 0.7039831095138017, "type": "data" }, { "id": "5d831813404ade3d94b8d377", "value": 0.6844200034628051, "type": "data" }, { "id": "5d831bed404ade3d94b8d61f", "value": 0.6767068069408032, "type": "data" }, { "id": "5d831cca404ade3d94b8d6e7", "value": 0.37322848358323635, "type": "data" }, { "id": "5d831d7c404ade3d94b8d7ab", "value": 0.3471557455239255, "type": "data" }, { "id": "5d8321fc404ade3d94b8d9ff", "value": 0.6113065286949759, "type": "data" }]
const thestat = { "ave": 0.5262098613410883, "cov": 0.24356357263129305, "max": 0.7282753159532952, "min": 0.32300643208654, "p95": 0.7068083185489441, "p99": 0.7226794287440242, "std": 0.1281655537820528 }
export default function ({data=thedata, stat=thestat}) {
    return (
        <Chart
            height={500}
            data={data}
            scale={{
                value: {
                    min: 0,
                    max: 1,
                    type: 'linear'
                },
                id: {
                    tickCount: data.length > 5 ? 5 : data.length
                }
            }}
            forceFit
        >
            <Legend />
            <Axis name='id' />
            <Axis name='value'/>
            <Tooltip/>
            <Geom type='line' position='id*value' shape='smooth' size={2} color={'type'} tooltip={['id*value*type', (x, y, type) => {
                return {
                    title: x,
                    name: type,
                    value: y
                }
            }]} />
        </Chart>
    )
}