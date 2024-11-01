Ext.define('PBS.D2DManagement', {
    extend: 'Ext.tab.Panel',
    alias: 'widget.pbsD2DManagement',

    title: 'Disk Backup',

    tools: [],

    border: true,
    defaults: {
	border: false,
	xtype: 'panel',
    },

    items: [
	{
	    xtype: 'pbsDiskBackupJobView',
	    title: gettext('Backup Jobs'),
	    itemId: 'd2d-backup-jobs',
	    iconCls: 'fa fa-floppy-o',
	},
	{
	    xtype: 'pbsDiskTargetPanel',
	    title: 'Targets',
	    itemId: 'targets',
	    iconCls: 'fa fa-desktop',
	},
    ],
});
