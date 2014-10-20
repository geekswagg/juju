from contextlib import contextmanager
from mock import patch
import os
import shutil
import tarfile
from tempfile import mkdtemp
from unittest import TestCase

from winbuildtest import (
    build_agent,
    build_client,
    create_cloud_agent,
    create_installer,
    GO_CMD,
    GOPATH,
    ISS_CMD,
)


@contextmanager
def temp_dir():
    dirname = mkdtemp()
    try:
        yield dirname
    finally:
        shutil.rmtree(dirname)


class WinBuildTestTestCase(TestCase):

    def test_build_client(self):
        with temp_dir() as cmd_dir:
            with temp_dir() as iss_dir:

                def make_juju(*args, **kwargs):
                    with open('%s/juju.exe' % cmd_dir, 'w') as fake_juju:
                        fake_juju.write('juju')

                with patch('winbuildtest.run',
                           return_value='', side_effect=make_juju) as run_mock:
                    build_client(cmd_dir, GO_CMD, GOPATH, iss_dir)
                    args, kwargs = run_mock.call_args
                    self.assertEqual((GO_CMD, 'build'), args)
                    self.assertEqual('386', kwargs['env'].get('GOARCH'))
                    self.assertEqual(GOPATH, kwargs['env'].get('GOPATH'))

    def test_create_installer(self):
        with temp_dir() as iss_dir:
            with temp_dir() as ci_dir:
                installer_name = 'juju-setup-1.20.1.exe'

                def make_installer(*args, **kwargs):
                    output_dir = os.path.join(iss_dir, 'output')
                    os.makedirs(output_dir)
                    installer_path = os.path.join(
                        output_dir, installer_name)
                    with open(installer_path, 'w') as fake_installer:
                        fake_installer.write('juju installer')

                with patch('winbuildtest.run',
                           return_value='',
                           side_effect=make_installer) as run_mock:
                    create_installer('1.20.1', iss_dir, ISS_CMD, ci_dir)
                    args, kwargs = run_mock.call_args
                    self.assertEqual((ISS_CMD, 'setup.iss'), args)
                    installer_path = os.path.join(ci_dir, installer_name)
                    self.assertTrue(os.path.isfile(installer_path))

    def test_build_agent(self):
        with temp_dir() as jujud_cmd_dir:
            with patch('winbuildtest.run', return_value='') as run_mock:
                build_agent(jujud_cmd_dir, GO_CMD, GOPATH)
                args, kwargs = run_mock.call_args
                self.assertEqual((GO_CMD, 'build'), args)
                self.assertEqual('amd64', kwargs['env'].get('GOARCH'))
                self.assertEqual(GOPATH, kwargs['env'].get('GOPATH'))

    def test_create_cloud_agent(self):
        with temp_dir() as cmd_dir:
            with temp_dir() as ci_dir:
                with open('%s/jujud.exe' % cmd_dir, 'w') as fake_jujud:
                    fake_jujud.write('jujud')
                create_cloud_agent('1.20.1', cmd_dir, ci_dir)
                agent = '{}/juju-1.20.1-win2012-amd64.tgz'.format(ci_dir)
                self.assertTrue(os.path.isfile(agent))
                with tarfile.open(name=agent, mode='r:gz') as tar:
                    self.assertEqual(['jujud.exe'], tar.getnames())
